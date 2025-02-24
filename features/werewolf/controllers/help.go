package controllers

import (
	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := `🎭 **Undercover Game Rules** 🎭

🃏 **Roles:**
- **Civilian:** Players who receive the same secret word.
- **Undercover:** Players who receive a slightly different word.

🏆 **Win Conditions:**
- **Civilians Win** by eliminating all Undercover.
- **Undercover Wins** when their number equals the Civilians.

📜 **Rules:**
1️⃣ Each player receives their role and corresponding word (or lack thereof).
2️⃣ Players take turns giving one-word clues about their word.
3️⃣ After each round, players vote to eliminate a suspect.
4️⃣ Eliminated players reveal their roles (by config)
   - Otherwise, the game continues until a win condition is met.

🔗 Get ready to bluff, deduce, and uncover the truth!

---

🎭 **Aturan Permainan Undercover** 🎭

🃏 **Peran:**
- **Civilian:** Pemain yang menerima kata rahasia yang sama.
- **Undercover:** Pemain yang menerima kata rahasia yang sedikit berbeda.

🏆 **Kondisi Menang:**
- **Civilian Menang** jika berhasil mengeliminasi semua Undercover.
- **Undercover Menang** jika jumlah mereka sama dengan Civilian.

📜 **Aturan:**
1️⃣ Setiap pemain menerima peran dan kata rahasia mereka.
2️⃣ Pemain bergiliran memberikan petunjuk satu kata tentang kata mereka.
3️⃣ Setelah setiap ronde, pemain memberikan suara untuk mengeliminasi seorang tersangka.
4️⃣ Pemain yang tereliminasi mengungkapkan perannya (berdasarkan konfigurasi).
   - Jika tidak, permainan berlanjut hingga salah satu kondisi menang terpenuhi.

🔗 Bersiaplah untuk menggertak, menganalisis, dan mengungkap kebenaran!`


	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
