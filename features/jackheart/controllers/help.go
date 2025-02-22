package controllers

import (
	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := `🎭 **Jack Heart Game Rules** 🎭

🔢 **Point System:**
- Starting Points: Total players × 2
- Maximum Points: Starting points + total players + 2

🃏 **Roles:**
- **Jack Heart**: The deceiver and main antagonist.
- **Pawn**: Neutral players working to identify Jack Heart.

🏆 **Win Conditions:**
- **Jack Heart Wins** by surviving until the end or reaching max points.
- **Pawn Win** by eliminating Jack Heart or reaching max points.

📜 **Rules:**
1️⃣ Symbols are hidden; players must spend 1 point to reveal another’s symbol.
2️⃣ Revealing a symbol costs -1 point, lying about a symbol earns +2 points, staying silent keeps points unchanged.
3️⃣ Players reaching 0 points are eliminated.
4️⃣ Voting occurs at the end of each round. The most voted player loses 3 points.

🔗 Get ready to deceive, deduce, and dominate!

---

🎭 **Aturan Permainan Jack Heart** 🎭

🔢 **Sistem Poin:**
- Poin Awal: Total pemain × 2
- Poin Maksimal: Poin awal + total pemain + 2

🃏 **Peran:**
- **Jack Heart**: Pembohong dan antagonis utama.
- **Pawn**: Pemain netral yang berusaha mengidentifikasi Jack Heart.

🏆 **Kondisi Menang:**
- **Jack Heart Menang** dengan bertahan hingga akhir atau mencapai poin maksimal.
- **Pawn Menang** dengan mengeliminasi Jack Heart atau mencapai poin maksimal.

📜 **Aturan:**
1️⃣ Simbol tersembunyi; pemain harus menggunakan 1 poin untuk mengungkap simbol pemain lain.
2️⃣ Mengungkap simbol mengurangi -1 poin, berbohong tentang simbol mendapatkan +2 poin, tetap diam poin tetap.
3️⃣ Pemain yang mencapai 0 poin akan tereliminasi.
4️⃣ Voting dilakukan di akhir setiap ronde. Pemain dengan suara terbanyak kehilangan 3 poin.

🔗 Bersiaplah untuk menipu, menganalisis, dan mendominasi!`
	

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral, 
		},
	})
}
