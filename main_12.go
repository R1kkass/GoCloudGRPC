package main

// import (
// 	"crypto/rand"
// 	"fmt"
// 	"math/big"
// )

// func main(){
//     GeneratePubKeys()
//     // fmt.Println(a, b)
// }

// func GeneratePubKeys() {
// 	max := new(big.Int)
// 	min := new(big.Int)

// 	var pp int64 = 22

// 	max = max.Exp(big.NewInt(2), big.NewInt(pp), nil).Sub(max, big.NewInt(1))
// 	min = min.Exp(big.NewInt(2), big.NewInt(pp-1), nil).Sub(min, big.NewInt(1))
// 	p, _ := rand.Int(rand.Reader, max.Sub(max, min).Add(max, min))
// 	g := generator(p);

// 	for !conditions(p, g) || !natural(p) {
// 		p, _ = rand.Int(rand.Reader, max.Sub(max, min).Add(max, min))

// 		g = generator(p);
// 	}

// 	max = max.Exp(big.NewInt(2), big.NewInt(pp), nil).Sub(max, big.NewInt(1))
// 	min = min.Exp(big.NewInt(2), big.NewInt(pp-1), nil).Sub(min, big.NewInt(1))
// 	a, _ := rand.Int(rand.Reader, max.Sub(p.Sub(p, big.NewInt(1)), min).Add(p.Sub(p, big.NewInt(1)), min))
// 	b, _ := rand.Int(rand.Reader, max.Sub(p.Sub(p, big.NewInt(1)), min).Add(p.Sub(p, big.NewInt(1)), min))

// 	s := make(chan *big.Int)
// 	s2 := make(chan *big.Int)

// 	go func () {
// 		var B = big.NewInt(0)
// 		B = B.Mod(B.Exp(g, b, nil), p)
// 		var s2 = big.NewInt(0)
// 		s2 = s2.Mod(s2.Exp(B, a, nil), p)
// 		s <- B
// 	}()

// 	go func () {
// 		var A = big.NewInt(0)
// 		A = A.Mod(A.Exp(g, a, nil), p)
// 		var s = big.NewInt(0)
// 		s = s.Mod(s.Exp(A, b, nil), p)
// 		s2 <-A
// 	}()
// 	fmt.Println(a,b,p)

// 	// var s3 = big.NewInt(0)
// 	// s3 = s3.Mod(s3.Exp(B, 10), nil, p)
// 	fmt.Println(<-s2,<-s)
// 	// return A, B
// }

// func natural(p *big.Int) bool {
// 	var n = big.NewInt(0)
//     for i := big.NewInt(2); i.Cmp(n.Sqrt(p)) <= 0 && i.Cmp(big.NewInt(40000))==-1; i = i.Add(i, big.NewInt(1)){
//         if n.Mod(p, i).Cmp(big.NewInt(0)) == 0 {
//             return false
//         }
//     }
//     return true
// }

// func powmod(a, b, p *big.Int) *big.Int {
//     res := big.NewInt(1)
//     for b.Cmp(big.NewInt(0)) == 1 {

//         if b.And(b, big.NewInt(1)).Cmp(big.NewInt(1)) == 0 {
//             res = res.Mul(res, a).Mod(res, p)
//             b = b.Sub(b, big.NewInt(1))
//         } else {
//             a = a.Mul(a, a).Mod(a, p)
//             b = b.Rsh(b, 1)
//         }
//     }
//     return res
// }

// func generator(p *big.Int) *big.Int {
// 	fact := []*big.Int{}
// 	phi := big.NewInt(0)
// 	phi = phi.Sub(p, big.NewInt(1))
// 	n := phi
//     var gn = big.NewInt(0)
// 	for i := big.NewInt(2); i.Mul(i, i).Cmp(n) == -1 || i.Mul(i, i).Cmp(n) == 0; i.Add(i, big.NewInt(1)) {
// 		if gn.Mod(n, i).Cmp(big.NewInt(0)) == 0 {
// 			fact = append(fact, i)
// 			for gn.Mod(n, i).Cmp(big.NewInt(0)) == 0 {

// 				n=n.Div(n, i)
// 			}
// 		}
// 	}
// 	if gn.Cmp(big.NewInt(1)) == +1 {
// 		fact = append(fact, n)
// 	}

// 	for res := big.NewInt(2); res.Cmp(p) != +1 && res.Cmp(big.NewInt(7)) == -1; res.Add(res, big.NewInt(1)) {
// 		ok := true
// 		for _, f := range fact {
// 			ok = ok && (powmod(res, gn.Div(phi,f), p).Cmp(big.NewInt(1)) != 0)
// 		}
// 		if ok {
// 			return res
// 		}
// 	}
// 	return big.NewInt(-1)
// }

// func conditions(p *big.Int, g *big.Int) bool {
// 	var n = big.NewInt(0)
// 	var gn = g

//     if n.Mod(p, big.NewInt(8)).Cmp(big.NewInt(7)) == 0 && gn.Cmp(big.NewInt(2)) == 0 {
//         return true
//     }
//     if n.Mod(p, big.NewInt(3)).Cmp(big.NewInt(2)) == 0 && gn.Cmp(big.NewInt(3)) == 0 {
//         return true
//     }
//     if gn.Cmp(big.NewInt(4)) == 0 {
//         return true
//     }
//     if (n.Mod(p, big.NewInt(5)).Cmp(big.NewInt(1)) == 0 || n.Mod(p, big.NewInt(5)).Cmp(big.NewInt(4)) == 0) && gn.Cmp(big.NewInt(5)) == 0 {
//         return true
//     }
//     if (n.Mod(p, big.NewInt(24)).Cmp(big.NewInt(19)) == 0 || n.Mod(p, big.NewInt(24)).Cmp(big.NewInt(23)) == 0) && gn.Cmp(big.NewInt(6)) == 0 {
//         return true
//     }
//     if (n.Mod(p, big.NewInt(7)).Cmp(big.NewInt(3)) == 0 || n.Mod(p, big.NewInt(7)).Cmp(big.NewInt(5)) == 0 || n.Mod(p, big.NewInt(7)).Cmp(big.NewInt(6)) == 0) && gn.Cmp(big.NewInt(7)) == 0{
//         return true
// 	}

// 	return false
// }