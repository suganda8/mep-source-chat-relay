use std::collections::hash_map::Entry;
use std::collections::HashMap;
use std::ffi::CStr;
use std::sync::Arc;

use lazy_static::lazy_static;

use tokio::runtime::Runtime;
use tokio::sync::RwLock;

use crate::packet::{Packet, Payload};
use crate::player::Player;
use crate::Result;

lazy_static! {
    static ref RUNTIME: Runtime = Runtime::new().unwrap();
}

/// Source Chat Relay client.
///
/// Measures need to be taken to ensure ptr lifetime on the shim.
pub struct Client(Arc<RwLock<ClientInner>>);

struct ClientInner {
    /// Map of players on the server.
    /// This will be updated upon client join/leave.
    players: HashMap<u64, Player>,
}

// By default *mut T is not safe for send.
// Need to ensure the shim side safety instead.
unsafe impl Send for Client {}

impl Default for Client {
    fn default() -> Self {
        Self(Arc::new(RwLock::new(ClientInner {
            players: HashMap::new(),
        })))
    }
}

impl Client {
    pub async fn receive_audio(&mut self, data: &[u8]) -> Result<()> {
        let mut packet = Packet::from_bytes(data)?;

        let header = packet.header()?;

        let payload = packet.payload()?;

        let mut inner = self.0.write().await;

        let player = match inner.players.entry(header.steam_id) {
            Entry::Occupied(p) => p.into_mut(),
            Entry::Vacant(v) => v.insert(Player::new()?),
        };

        match payload {
            Payload::OpusPLC(data) => {
                println!("!!!!! Opus PLC {}", data.len());

                match player.transcode(data) {
                    Ok(mut d) => {
                        println!("ok transcode {}", d.len());
                    }
                    Err(e) => println!("{:?}", e),
                }
            }
            Payload::Silence(ns) => {
                println!("!!!!! Silence {}", ns);

                // Silence payload should also be sent on the wire.
            }
        }

        Ok(())
    }

    pub async fn client_put_in_server(&mut self, steamid: u64, name: &str) -> Result<()> {
        let mut inner = self.0.write().await;

        inner.players.insert(steamid, Player::new()?);

        Ok(())
    }

    pub async fn client_disconnect(&mut self, steamid: u64) -> Result<()> {
        let mut inner = self.0.write().await;

        inner.players.remove(&steamid);

        Ok(())
    }
}

#[no_mangle]
pub extern "C" fn new_client() -> *mut Client {
    let c = Client::default();
    let b = Box::new(c);

    Box::into_raw(b)
}

#[no_mangle]
pub unsafe extern "C" fn receive_audio(
    client: *mut Client,
    bytes: i32,
    data: *const std::os::raw::c_char,
) {
    if !client.is_null() {
        let d = std::slice::from_raw_parts(data as *const u8, bytes as usize).to_owned();

        let client = &mut *client;

        RUNTIME.spawn(async move {
            let _ = client.receive_audio(&d).await;
        });
    }
}

#[no_mangle]
pub unsafe extern "C" fn client_put_in_server(
    client: *mut Client,
    steamid: u64,
    name: *const std::os::raw::c_char,
) {
    if !client.is_null() {
        let name = CStr::from_ptr(name).to_string_lossy();

        let client = &mut *client;

        RUNTIME.spawn(async move {
            let _ = client.client_put_in_server(steamid, &name).await;
        });
    }
}

#[no_mangle]
pub unsafe extern "C" fn client_disconnect(client: *mut Client, steamid: u64) {
    if !client.is_null() {
        let client = &mut *client;

        RUNTIME.spawn(async move {
            let _ = client.client_disconnect(steamid).await;
        });
    }
}

#[no_mangle]
pub unsafe extern "C" fn free_client(client: *mut Client) {
    if client.is_null() {
        return;
    }

    Box::from_raw(client);
}