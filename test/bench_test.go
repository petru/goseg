package goseg_test

import (
	"testing"

	"goseg"
)

const benchText = `Hello World. My name is Jonas. What is your name? My name is Jonas. There it is! I found it. My name is Jonas E. Smith. Please turn to p. 55. Were Jane and co. at the party? They closed the deal with Pitt, Briggs & Co. at noon. Let's ask Jane and co. They should know. They closed the deal with Pitt, Briggs & Co. It closed yesterday. I can see Mt. Fuji from here. St. Michael's Church is on 5th st. near the light. That is JFK Jr.'s book. I visited the U.S.A. last year. I live in the E.U. How about you? I live in the U.S. How about you? I work for the U.S. Government in Virginia. I have lived in the U.S. for 20 years. She has $100.00 in her bag. She has $100.00. It is in her bag. He teaches science (He previously worked for 5 years as an engineer.) at the local University. Her email is Jane.Doe@example.com. I sent her an email. The site is: https://www.example.50.com/new-site/awesome_content.html. Please check it out. She turned to him, 'This is great.' she said. She turned to him, "This is great." she said. She turned to him, "This is great." She held the book out to show him. Hello!! Long time no see. Hello?? Who is there? Hello!? Is that you? Hello?! Is that you? 1.) The first item 2.) The second item 1.) The first item. 2.) The second item. 1) The first item 2) The second item 1) The first item. 2) The second item. 1. The first item 2. The second item 1. The first item. 2. The second item. • 9. The first item • 10. The second item ⁃9. The first item ⁃10. The second item a. The first item b. The second item c. The third list item This is a sentence
cut off in the middle because pdf. It was a cold 
night in the city. features
contact manager
events, activities
 You can find it at N°. 1026.253.553. That is where the treasure is. She works at Yahoo! in the accounting department. We make a good team, you and I. Did you see Albert I. Jones yesterday? Thoreau argues that by simplifying one’s life, “the laws of the universe will appear less complex. . . .” "Bohr [...] used the analogy of parallel stairways [...]" (Smith 55). If words are left off at the end of a sentence, and that is all that is omitted, indicate the omission with ellipsis marks (preceded and followed by a space) and then indicate the end of the sentence with a period . . . . Next sentence. I never meant that.... She left the store. I wasn’t really ... well, what I mean...see . . . what I'm saying, the thing is . . . I didn’t mean it. One further habit which was somewhat weakened . . . was that of combining words into self-interpreting compounds. . . . The practice was not abandoned. . . .`

func BenchmarkSegmentGoldenRulesx100(b *testing.B) {
	text := ""
	for i := 0; i < 100; i++ {
		text += benchText
	}
	_ = goseg.Segment(text[:len(benchText)], goseg.Language("en"))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = goseg.Segment(text, goseg.Language("en"))
	}
}

func BenchmarkSegmentGoldenRulesOnce(b *testing.B) {
	_ = goseg.Segment(benchText, goseg.Language("en"))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = goseg.Segment(benchText, goseg.Language("en"))
	}
}
