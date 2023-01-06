// To parse and unparse this JSON data :
// emojis, err := UnmarshalEmojis(bytes)
// bytes, err = emojis.Marshal()

package emoji

import (
	"os"
	"regexp"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Emojis []EmojisElement

func UnmarshalEmojis(data []byte) (Emojis, error) {
	var r Emojis

	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Emojis) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type EmojisElement struct {
	Emoji          string         `json:"emoji"`
	Description    string         `json:"description"`
	Category       Category       `json:"category"`
	Aliases        []string       `json:"aliases"`
	Tags           []string       `json:"tags"`
	UnicodeVersion UnicodeVersion `json:"unicode_version"`
	IosVersion     string         `json:"ios_version"`
	SkinTones      *bool          `json:"skin_tones,omitempty"`
}

type Category string

const (
	Activities     Category = "Activities"
	AnimalsNature  Category = "Animals & Nature"
	Flags          Category = "Flags"
	FoodDrink      Category = "Food & Drink"
	Objects        Category = "Objects"
	PeopleBody     Category = "People & Body"
	SmileysEmotion Category = "Smileys & Emotion"
	Symbols        Category = "Symbols"
	TravelPlaces   Category = "Travel & Places"
)

type UnicodeVersion string

const (
	Empty  UnicodeVersion = ""
	The110 UnicodeVersion = "11.0"
	The120 UnicodeVersion = "12.0"
	The121 UnicodeVersion = "12.1"
	The130 UnicodeVersion = "13.0"
	The131 UnicodeVersion = "13.1"
	The140 UnicodeVersion = "14.0"
	The30  UnicodeVersion = "3.0"
	The32  UnicodeVersion = "3.2"
	The40  UnicodeVersion = "4.0"
	The41  UnicodeVersion = "4.1"
	The51  UnicodeVersion = "5.1"
	The52  UnicodeVersion = "5.2"
	The60  UnicodeVersion = "6.0"
	The61  UnicodeVersion = "6.1"
	The70  UnicodeVersion = "7.0"
	The80  UnicodeVersion = "8.0"
	The90  UnicodeVersion = "9.0"
)

func DecodeEmojisToAliases(str string) (res string, err error) {
	// os.Open require defering os.File.Close() and then io.ReadAll() to read the bytes. This behavior because opening file may not only to read but also write.
	// On the other side, os.ReadFile only read file and pass the bytes directly.
	emojiJson, errJson := os.ReadFile("../server/emoji/emoji.json")

	// Reminder, if use the string directly to MustCompile, use "Backtick" (``) instead normal Double-quotes.
	// regexp.MustCompile(`long-emoji-regex`)
	emojiRegex, errRegex := os.ReadFile("../server/emoji/emoji-regex.txt")

	if errJson != nil || errRegex != nil {
		return "", err
	}

	emojis, err := UnmarshalEmojis(emojiJson)

	if err != nil {
		return "", err
	}

	emojiRegexStr := string(emojiRegex)

	re := regexp.MustCompile(emojiRegexStr)

	replaced := re.ReplaceAllStringFunc(str, func(match string) string {
		for i := 0; i < len(emojis); i++ {
			if match == emojis[i].Emoji {
				if len(emojis[i].Aliases) > 0 {
					return ":" + emojis[i].Aliases[0] + ":"
				}
			}
		}
		return ""
	})

	return replaced, nil
}
