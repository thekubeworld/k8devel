package emoji

/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"strconv"
	"strings"
)

type Type struct {
	Rocket string              // 🚀
	ThumbsUp string            // 👍
	ThumbsDown string          // 👎
	Collision string           // 💥
	Robot string               // 🤖
	AlienMonster string        // 👾 
	Alien string               // 👽
	Skull string               // 💀
	SkullAndCrossBones string  // ☠
	Ghost string               // 👻
	AngryFace string           // 😠
	NerdFace string            // 🤓
	SmileFace string           // 😀  ￼
	PartyFace string	   // 🥳 
	SatelliteAntenna string    // 📡
	CheckMarkButton string     // ✅
	CrossMark string           // ❌
	ChequeredFlag string	   // 🏁
	MegaPhone string	   // 📣
        Rainbow string		   // 🌈
        HourGlassNotDone string    // ⏳
	StopSign string		   // 🛑
	Construction string	   // 🚧
	RedHeart string		   // ❤
	PileOfPoo string	   // 💩
	ClownFace string	   // 🤡
	SleepingFace string	   // 😴
}

// LoadEmojis will load emojis codes
//
//	Table of codes:
//	http://www.unicode.org/emoji/charts/emoji-list.html
//	https://www.unicode.org/emoji/charts/full-emoji-modifiers.html
//
// Args:
//     
//      None
//
// Returns:
//
//	Emoji Type struct
//
func LoadEmojis() Type {

	Emojis := Type {
		Rocket: "\\U0001f680",            // 🚀
		ThumbsUp: "\\U0001f44D",          // 👍
		ThumbsDown: "\\U0001f44e",        // 👎
		Collision: "\\U0001f4a5",         // 💥
		Robot: "\\U0001f916",             // 🤖
		AlienMonster: "\\U0001f47e",      // 👾
		Alien: "\\U0001f47d",             // 👽
		Skull: "\\U0001f480",             // 💀
		SkullAndCrossBones: "\\U0002620", // ☠
		Ghost: "\\U0001f47b",             // 👻
		AngryFace: "\\U0001f620",         // 😠
		NerdFace: "\\U0001f913",          // 🤓
		SmileFace: "\\U0001f60e",         // 😀
		PartyFace: "\\U0001f973",         // 🥳 
		SatelliteAntenna: "\\U0001f4e1",  // 📡
		CheckMarkButton: "\\U0002705",    // ✅
		CrossMark: "\\U000274C",          // ❌
		ChequeredFlag: "\\U0001f3c1",	  // 🏁
		MegaPhone: "\\U0001f4e3",	  // 📣
		Rainbow: "\\U0001f308",           // 🌈
		HourGlassNotDone: "\\U00023f3",   // ⏳
		StopSign: "\\U0001f6D1",	  // 🛑
		Construction: "\\U0001f6a7",	  // 🚧
		RedHeart: "\\U0002764",           // ❤
		PileOfPoo: "\\U0001f4a9",	  // 💩
		ClownFace: "\\U0001f921",	  // 🤡
		SleepingFace: "\\U0001f634",	  // 😴
	}
	return Emojis
}

// UnquoteCode will unquote code for the
// emoji
//
//	Table of codes:
//	http://www.unicode.org/emoji/charts/emoji-list.html
//	https://www.unicode.org/emoji/charts/full-emoji-modifiers.html
//
// Args:
//     
//      string - unicode string quoted
//
// Returns:
//
//	string or error
//
func Show(s string) string {
	ret, _ := strconv.ParseInt(
		strings.TrimPrefix(s, "\\U"), 16, 32)

	return string(ret)
}
