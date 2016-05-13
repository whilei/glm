package geo2

import (
	"github.com/luxengine/glm"
	"github.com/luxengine/glm/geo2/internal/qhull"
	"github.com/luxengine/math"
	"testing"
)

func TestQuickhull(t *testing.T) {
	// These tests we're visualised and verified by hand (using an svg debug
	// program), then the data was exported here. If the algorithm is changed
	// and this breaks it may not be that the algorithm is broken it may just be
	// that the data is now in a different order (this test only does slice
	// content equality).
	tests := []struct {
		points []glm.Vec2
		hull   []glm.Vec2
	}{
		{
			points: []glm.Vec2{{0, 0}, {0.5, -0.1}, {1, 0}, {1.1, 1}, {0, 1}, {0.5, 0.5}, {-0.4, 0.5}, {0.4, 0.4}, {0.5, 0.4}, {0.6, 0.6}},
			hull:   []glm.Vec2{{-0.4, 0.5}, {0, 1}, {1.1, 1}, {1, 0}, {0.5, -0.1}, {0, 0}},
		},
		{
			points: []glm.Vec2{{0.5, 0.5}, {-0.4, 0.5}, {0.4, 0.4}, {0.5, 0.4}, {0.6, 0.6}, {0, 0}, {0.5, -0.1}, {1, 0}, {1.1, 1}, {0, 1}},
			hull:   []glm.Vec2{{-0.4, 0.5}, {0, 1}, {1.1, 1}, {1, 0}, {0.5, -0.1}, {0, 0}},
		},
		{
			points: []glm.Vec2{{270, 319}, {373, 174}, {396, 92}, {321, 148}, {354, 416}, {290, 322}, {9, 457}, {361, 328}, {421, 49}, {486, 256}},
			hull:   []glm.Vec2{{9, 457}, {354, 416}, {486, 256}, {421, 49}, {321, 148}},
		},
		{
			points: []glm.Vec2{{62, 482}, {448, 208}, {393, 395}, {75, 180}, {201, 136}, {339, 22}, {354, 491}, {7, 369}, {155, 108}, {252, 352}},
			hull:   []glm.Vec2{{7, 369}, {62, 482}, {354, 491}, {393, 395}, {448, 208}, {339, 22}, {155, 108}, {75, 180}},
		},
		{
			points: []glm.Vec2{{378, 112}, {239, 65}, {202, 239}, {277, 240}, {451, 269}, {498, 127}, {322, 490}, {115, 442}, {415, 119}, {135, 496}},
			hull:   []glm.Vec2{{115, 442}, {135, 496}, {322, 490}, {451, 269}, {498, 127}, {239, 65}},
		},
		{
			points: []glm.Vec2{{43, 429}, {136, 262}, {79, 17}, {383, 374}, {216, 39}, {61, 347}, {244, 343}, {308, 172}, {283, 422}, {405, 252}},
			hull:   []glm.Vec2{{43, 429}, {283, 422}, {383, 374}, {405, 252}, {216, 39}, {79, 17}},
		},
		{
			points: []glm.Vec2{{65, 427}, {150, 399}, {242, 55}, {99, 136}, {182, 339}, {344, 185}, {38, 109}, {410, 332}, {366, 41}, {133, 254}},
			hull:   []glm.Vec2{{38, 109}, {65, 427}, {410, 332}, {366, 41}, {242, 55}},
		},
		{
			points: []glm.Vec2{{379, 227}, {111, 451}, {134, 225}, {3, 441}, {163, 196}, {296, 452}, {356, 46}, {63, 196}, {73, 410}, {448, 214}},
			hull:   []glm.Vec2{{3, 441}, {111, 451}, {296, 452}, {448, 214}, {356, 46}, {63, 196}},
		},
		{
			points: []glm.Vec2{{33, 349}, {441, 345}, {7, 98}, {82, 46}, {283, 303}, {184, 39}, {179, 210}, {433, 58}, {81, 376}, {157, 492}},
			hull:   []glm.Vec2{{7, 98}, {33, 349}, {157, 492}, {441, 345}, {433, 58}, {184, 39}, {82, 46}},
		},
		{
			points: []glm.Vec2{{12, 60}, {235, 399}, {344, 439}, {215, 225}, {424, 141}, {456, 471}, {129, 421}, {42, 31}, {483, 490}, {16, 90}},
			hull:   []glm.Vec2{{12, 60}, {16, 90}, {129, 421}, {483, 490}, {424, 141}, {42, 31}},
		},
		{
			points: []glm.Vec2{{178, 453}, {47, 375}, {10, 147}, {240, 111}, {197, 135}, {326, 258}, {19, 44}, {449, 350}, {439, 201}, {454, 409}},
			hull:   []glm.Vec2{{10, 147}, {47, 375}, {178, 453}, {454, 409}, {439, 201}, {240, 111}, {19, 44}},
		},
		{
			points: []glm.Vec2{{90, 352}, {243, 287}, {368, 105}, {55, 476}, {427, 44}, {386, 288}, {303, 100}, {4, 47}, {223, 172}, {352, 267}},
			hull:   []glm.Vec2{{4, 47}, {55, 476}, {386, 288}, {427, 44}},
		},
		{
			points: []glm.Vec2{{165, 27}, {495, 340}, {472, 28}, {403, 158}, {300, 277}, {107, 231}, {471, 462}, {445, 38}, {266, 495}, {487, 180}},
			hull:   []glm.Vec2{{107, 231}, {266, 495}, {471, 462}, {495, 340}, {487, 180}, {472, 28}, {165, 27}},
		},
		{
			points: []glm.Vec2{{469, 72}, {173, 361}, {364, 48}, {356, 466}, {274, 199}, {252, 418}, {78, 309}, {424, 181}, {265, 45}, {493, 162}},
			hull:   []glm.Vec2{{78, 309}, {252, 418}, {356, 466}, {493, 162}, {469, 72}, {364, 48}, {265, 45}},
		},
		{
			points: []glm.Vec2{{105, 414}, {26, 328}, {409, 182}, {314, 144}, {364, 30}, {0, 180}, {424, 112}, {351, 19}, {409, 353}, {2, 223}},
			hull:   []glm.Vec2{{0, 180}, {2, 223}, {26, 328}, {105, 414}, {409, 353}, {424, 112}, {364, 30}, {351, 19}},
		},
		{
			points: []glm.Vec2{{217, 311}, {130, 34}, {24, 381}, {450, 78}, {413, 409}, {495, 231}, {440, 109}, {127, 282}, {325, 468}, {20, 59}},
			hull:   []glm.Vec2{{20, 59}, {24, 381}, {325, 468}, {413, 409}, {495, 231}, {450, 78}, {130, 34}},
		},
		{
			points: []glm.Vec2{{463, 367}, {350, 449}, {188, 476}, {163, 67}, {264, 323}, {54, 129}, {189, 265}, {112, 195}, {115, 471}, {40, 492}},
			hull:   []glm.Vec2{{40, 492}, {188, 476}, {350, 449}, {463, 367}, {163, 67}, {54, 129}},
		},
		{
			points: []glm.Vec2{{168, 343}, {228, 411}, {471, 430}, {95, 473}, {198, 9}, {265, 194}, {330, 408}, {214, 309}, {270, 261}, {199, 60}},
			hull:   []glm.Vec2{{95, 473}, {471, 430}, {198, 9}},
		},
		{
			points: []glm.Vec2{{369, 122}, {181, 433}, {383, 388}, {45, 37}, {396, 201}, {466, 22}, {479, 154}, {437, 168}, {23, 462}, {373, 351}},
			hull:   []glm.Vec2{{23, 462}, {181, 433}, {383, 388}, {479, 154}, {466, 22}, {45, 37}},
		},
		{
			points: []glm.Vec2{{270, 176}, {437, 60}, {14, 438}, {473, 487}, {355, 147}, {407, 152}, {254, 396}, {3, 294}, {298, 375}, {493, 220}},
			hull:   []glm.Vec2{{3, 294}, {14, 438}, {473, 487}, {493, 220}, {437, 60}},
		},
		{
			points: []glm.Vec2{{146, 389}, {415, 404}, {51, 325}, {428, 298}, {365, 441}, {183, 201}, {334, 235}, {473, 374}, {359, 261}, {102, 405}},
			hull:   []glm.Vec2{{51, 325}, {102, 405}, {365, 441}, {473, 374}, {428, 298}, {334, 235}, {183, 201}},
		},
		{
			points: []glm.Vec2{{448, 395}, {195, 479}, {353, 300}, {276, 496}, {102, 110}, {108, 233}, {420, 329}, {324, 82}, {227, 423}, {126, 49}},
			hull:   []glm.Vec2{{102, 110}, {108, 233}, {195, 479}, {276, 496}, {448, 395}, {324, 82}, {126, 49}},
		},
		{
			points: []glm.Vec2{{22, 61}, {328, 449}, {289, 59}, {35, 82}, {134, 40}, {224, 490}, {100, 319}, {400, 146}, {165, 133}, {47, 145}},
			hull:   []glm.Vec2{{22, 61}, {47, 145}, {100, 319}, {224, 490}, {328, 449}, {400, 146}, {289, 59}, {134, 40}},
		},
		{
			points: []glm.Vec2{{163, 90}, {40, 319}, {285, 289}, {402, 328}, {237, 114}, {438, 218}, {65, 146}, {287, 115}, {0, 357}, {362, 466}},
			hull:   []glm.Vec2{{0, 357}, {362, 466}, {438, 218}, {287, 115}, {163, 90}, {65, 146}},
		},
		{
			points: []glm.Vec2{{323, 76}, {279, 49}, {9, 300}, {133, 166}, {454, 451}, {422, 225}, {429, 35}, {427, 348}, {404, 477}, {492, 153}},
			hull:   []glm.Vec2{{9, 300}, {404, 477}, {454, 451}, {492, 153}, {429, 35}, {279, 49}, {133, 166}},
		},
		{
			points: []glm.Vec2{{456, 199}, {193, 409}, {53, 446}, {37, 71}, {155, 137}, {44, 83}, {424, 356}, {328, 292}, {219, 169}, {106, 449}},
			hull:   []glm.Vec2{{37, 71}, {53, 446}, {106, 449}, {424, 356}, {456, 199}},
		},
		{
			points: []glm.Vec2{{171, 438}, {256, 371}, {453, 118}, {221, 397}, {232, 203}, {110, 345}, {237, 74}, {77, 17}, {351, 430}, {421, 204}},
			hull:   []glm.Vec2{{77, 17}, {110, 345}, {171, 438}, {351, 430}, {453, 118}},
		},
		{
			points: []glm.Vec2{{462, 225}, {427, 212}, {462, 34}, {79, 453}, {389, 7}, {247, 405}, {242, 160}, {342, 323}, {69, 164}, {205, 199}},
			hull:   []glm.Vec2{{69, 164}, {79, 453}, {247, 405}, {462, 225}, {462, 34}, {389, 7}},
		},
		{
			points: []glm.Vec2{{307, 332}, {412, 83}, {23, 426}, {80, 103}, {126, 27}, {147, 130}, {442, 292}, {127, 93}, {153, 238}, {485, 398}},
			hull:   []glm.Vec2{{23, 426}, {485, 398}, {412, 83}, {126, 27}, {80, 103}},
		},
		{
			points: []glm.Vec2{{123, 26}, {197, 334}, {170, 349}, {405, 354}, {132, 359}, {373, 356}, {209, 3}, {54, 84}, {401, 401}, {86, 398}},
			hull:   []glm.Vec2{{54, 84}, {86, 398}, {401, 401}, {405, 354}, {209, 3}, {123, 26}},
		},
		{
			points: []glm.Vec2{{42, 257}, {68, 185}, {370, 469}, {195, 153}, {320, 401}, {483, 302}, {473, 441}, {494, 457}, {379, 89}, {12, 92}},
			hull:   []glm.Vec2{{12, 92}, {42, 257}, {370, 469}, {494, 457}, {483, 302}, {379, 89}},
		},
		{
			points: []glm.Vec2{{351, 81}, {304, 171}, {466, 139}, {476, 289}, {8, 409}, {457, 295}, {453, 136}, {116, 494}, {70, 273}, {314, 144}},
			hull:   []glm.Vec2{{8, 409}, {116, 494}, {476, 289}, {466, 139}, {351, 81}, {70, 273}},
		},
		{
			points: []glm.Vec2{{142, 55}, {56, 427}, {439, 65}, {320, 23}, {310, 464}, {269, 445}, {402, 446}, {273, 61}, {76, 172}, {345, 225}},
			hull:   []glm.Vec2{{56, 427}, {310, 464}, {402, 446}, {439, 65}, {320, 23}, {142, 55}, {76, 172}},
		},
		{
			points: []glm.Vec2{{115, 422}, {91, 125}, {2, 249}, {412, 66}, {276, 247}, {83, 386}, {115, 221}, {70, 117}, {40, 390}, {348, 282}},
			hull:   []glm.Vec2{{2, 249}, {40, 390}, {115, 422}, {348, 282}, {412, 66}, {70, 117}},
		},
		{
			points: []glm.Vec2{{289, 372}, {310, 338}, {310, 86}, {88, 352}, {486, 99}, {327, 304}, {144, 226}, {135, 146}, {249, 352}, {427, 99}},
			hull:   []glm.Vec2{{88, 352}, {289, 372}, {486, 99}, {310, 86}, {135, 146}},
		},
		{
			points: []glm.Vec2{{486, 217}, {52, 24}, {352, 56}, {326, 368}, {11, 363}, {186, 205}, {67, 406}, {435, 45}, {488, 350}, {406, 279}},
			hull:   []glm.Vec2{{11, 363}, {67, 406}, {488, 350}, {486, 217}, {435, 45}, {52, 24}},
		},
		{
			points: []glm.Vec2{{426, 443}, {213, 7}, {412, 5}, {71, 466}, {73, 315}, {33, 290}, {201, 263}, {227, 222}, {159, 129}, {475, 441}},
			hull:   []glm.Vec2{{33, 290}, {71, 466}, {475, 441}, {412, 5}, {213, 7}},
		},
		{
			points: []glm.Vec2{{437, 148}, {446, 304}, {260, 198}, {401, 167}, {342, 111}, {308, 497}, {192, 326}, {17, 443}, {245, 305}, {253, 490}},
			hull:   []glm.Vec2{{17, 443}, {253, 490}, {308, 497}, {446, 304}, {437, 148}, {342, 111}},
		},
		{
			points: []glm.Vec2{{213, 18}, {147, 496}, {136, 377}, {250, 461}, {197, 346}, {322, 0}, {251, 39}, {246, 476}, {212, 238}, {485, 217}},
			hull:   []glm.Vec2{{136, 377}, {147, 496}, {246, 476}, {485, 217}, {322, 0}, {213, 18}},
		},
		{
			points: []glm.Vec2{{282, 88}, {196, 126}, {460, 81}, {348, 341}, {353, 454}, {69, 451}, {388, 386}, {164, 180}, {13, 410}, {89, 350}},
			hull:   []glm.Vec2{{13, 410}, {69, 451}, {353, 454}, {388, 386}, {460, 81}, {282, 88}, {196, 126}},
		},
		{
			points: []glm.Vec2{{318, 274}, {235, 93}, {335, 253}, {164, 256}, {471, 398}, {306, 432}, {352, 202}, {336, 124}, {348, 289}, {292, 99}},
			hull:   []glm.Vec2{{164, 256}, {306, 432}, {471, 398}, {336, 124}, {292, 99}, {235, 93}},
		},
		{
			points: []glm.Vec2{{498, 338}, {358, 187}, {107, 196}, {70, 433}, {175, 313}, {452, 73}, {305, 363}, {308, 339}, {128, 409}, {214, 383}},
			hull:   []glm.Vec2{{70, 433}, {498, 338}, {452, 73}, {107, 196}},
		},
		{
			points: []glm.Vec2{{204, 131}, {74, 370}, {303, 294}, {59, 292}, {121, 149}, {78, 457}, {427, 26}, {184, 497}, {111, 277}, {80, 377}},
			hull:   []glm.Vec2{{59, 292}, {78, 457}, {184, 497}, {303, 294}, {427, 26}, {121, 149}},
		},
		{
			points: []glm.Vec2{{317, 77}, {397, 139}, {352, 197}, {183, 23}, {156, 75}, {262, 451}, {446, 399}, {224, 343}, {385, 122}, {388, 273}},
			hull:   []glm.Vec2{{156, 75}, {224, 343}, {262, 451}, {446, 399}, {397, 139}, {385, 122}, {317, 77}, {183, 23}},
		},
		{
			points: []glm.Vec2{{463, 245}, {83, 249}, {347, 466}, {103, 100}, {372, 405}, {417, 421}, {274, 166}, {226, 169}, {202, 467}, {320, 401}},
			hull:   []glm.Vec2{{83, 249}, {202, 467}, {347, 466}, {417, 421}, {463, 245}, {274, 166}, {103, 100}},
		},
		{
			points: []glm.Vec2{{401, 359}, {429, 366}, {148, 455}, {360, 334}, {306, 200}, {308, 339}, {26, 93}, {133, 470}, {336, 435}, {296, 260}},
			hull:   []glm.Vec2{{26, 93}, {133, 470}, {336, 435}, {429, 366}, {306, 200}},
		},
		{
			points: []glm.Vec2{{133, 438}, {247, 95}, {413, 423}, {5, 237}, {429, 427}, {121, 33}, {201, 279}, {490, 355}, {449, 56}, {378, 102}},
			hull:   []glm.Vec2{{5, 237}, {133, 438}, {429, 427}, {490, 355}, {449, 56}, {121, 33}},
		},
		{
			points: []glm.Vec2{{331, 389}, {442, 199}, {357, 155}, {213, 238}, {460, 269}, {302, 320}, {213, 385}, {346, 9}, {347, 75}, {281, 419}},
			hull:   []glm.Vec2{{213, 238}, {213, 385}, {281, 419}, {331, 389}, {460, 269}, {442, 199}, {346, 9}},
		},
		{
			points: []glm.Vec2{{261, 201}, {408, 98}, {66, 343}, {287, 5}, {269, 430}, {285, 311}, {487, 32}, {255, 448}, {225, 318}, {78, 305}},
			hull:   []glm.Vec2{{66, 343}, {255, 448}, {269, 430}, {487, 32}, {287, 5}, {78, 305}},
		},
		{
			points: []glm.Vec2{{22, 449}, {31, 210}, {168, 440}, {257, 206}, {221, 426}, {411, 421}, {218, 399}, {431, 68}, {154, 319}, {158, 326}},
			hull:   []glm.Vec2{{22, 449}, {168, 440}, {411, 421}, {431, 68}, {31, 210}},
		},
		{
			points: []glm.Vec2{{362, 76}, {110, 332}, {166, 194}, {483, 373}, {111, 74}, {34, 26}, {267, 366}, {446, 311}, {245, 443}, {400, 235}},
			hull:   []glm.Vec2{{34, 26}, {110, 332}, {245, 443}, {483, 373}, {362, 76}},
		},
		{
			points: []glm.Vec2{{403, 316}, {486, 168}, {259, 294}, {231, 466}, {338, 440}, {285, 106}, {290, 490}, {15, 249}, {106, 244}, {110, 171}},
			hull:   []glm.Vec2{{15, 249}, {231, 466}, {290, 490}, {338, 440}, {486, 168}, {285, 106}, {110, 171}},
		},
		{
			points: []glm.Vec2{{290, 21}, {428, 49}, {424, 377}, {88, 268}, {37, 41}, {187, 275}, {276, 317}, {170, 330}, {166, 184}, {72, 449}},
			hull:   []glm.Vec2{{37, 41}, {72, 449}, {424, 377}, {428, 49}, {290, 21}},
		},
	}

	for i, test := range tests {
		hull := Quickhull(test.points)
		for n, v := range hull.Vertices {
			if !v.Position.ApproxEqual(&test.hull[n]) {
				t.Errorf("[%d] %+v", i, test.points)
				t.Errorf("[%d]\n\thull %v\n\twant %v", i, hull, test.hull)
				break
			}
		}
	}
}

func TestQuickhullSupport(t *testing.T) {
	points := []glm.Vec2{{-0.1, -0.1}, {0.5, -0.1}, {1, 0}, {1.1, 1}, {0, 1}, {0.5, 0.5}, {-0.4, 0.5}, {0.4, 0.4}, {0.5, 0.4}, {0.6, 0.6}}
	hull := Quickhull(points)

	var pos []glm.Vec2
	for n := range hull.Vertices {
		pos = append(pos, hull.Vertices[n].Position)
	}
	t.Log("Vertices ", pos)
	t.Log("support dir ", qhull.SupportDirection)
	t.Log("support cache ", hull.Vertices[hull.bestSupport[0]].Position,
		hull.Vertices[hull.bestSupport[1]].Position,
		hull.Vertices[hull.bestSupport[2]].Position)

	const sep = 7
	for n := 0; n < sep; n++ {
		dir := glm.Vec2{math.Cos(float32(n) * 2.0 * math.Pi / float32(sep)), math.Sin(float32(n) * 2.0 * math.Pi / float32(sep))}
		s := hull.Support(&dir)
		t.Logf("[%d] %s %d", n, dir.String(), s)
	}
}
