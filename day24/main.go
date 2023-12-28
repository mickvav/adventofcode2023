package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/pa-m/optimize"
	"gonum.org/v1/gonum/mat"
)

type Input struct {
	m []string
	h []Hail
}

type pos struct {
	x, y, z float64
}

type dir pos

type Hail struct {
	p pos
	d dir
}

func ReadHail(s string) Hail {
	p := strings.Split(s, " @ ")
	h := Hail{}
	ps := strings.Split(p[0], ", ")
	ds := strings.Split(p[1], ", ")
	h.p.x, _ = strconv.ParseFloat(strings.Trim(ps[0], " "), 64)
	h.p.y, _ = strconv.ParseFloat(strings.Trim(ps[1], " "), 64)
	h.p.z, _ = strconv.ParseFloat(strings.Trim(ps[2], " "), 64)
	h.d.x, _ = strconv.ParseFloat(strings.Trim(ds[0], " "), 64)
	h.d.y, _ = strconv.ParseFloat(strings.Trim(ds[1], " "), 64)
	h.d.z, _ = strconv.ParseFloat(strings.Trim(ds[2], " "), 64)
	return h
}
func ReadInput(scanner *bufio.Scanner) Input {
	res := Input{
		m: []string{},
	}
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		res.m = append(res.m, line)
		res.h = append(res.h, ReadHail(line))
		row++
	}
	return res
}

var minCoord = float64(200000000000000)
var maxCoord = float64(400000000000000)

func (h1 Hail) Intersect2d(h2 Hail) (pos, error) {

	D := h1.d.y*h2.d.x - h1.d.x*h2.d.y
	if D == 0 {
		return pos{}, errors.New("parallel")
	}
	// Note - no 0 speeds!
	// x1(t) = x10 + vx1*t => t=(x1-x10)/vx1
	// y1(t) = y10 + vy1*t
	//
	// y1(x1) = y10 + vy1*(x1-x10)/vx1
	// y2(x2) = y20 + vy2*(x2-x20)/vx2
	//
	// Intersection:
	// y10 + vy1*(x-x10)/vx1 == y20 + vy2*(x-x20)/vx2
	// solve for x
	// x*(vy1/vx1 - vy2/vx2) == y20 - y10 + vy1*x10/vx1 - vy2*x20/vx2
	//
	// x = ((y20 - y10)*(vx1*vx2) + vy1 * x10 * vx2 - vy2 * x20 * vx1 ) / (vy1 * vx2 - vx1 * vy2)

	// (symmetry:)
	// y = ((x20 - x10)*(vy1*vy2) + vx1 * y10 * vy2 - vx2 * y20 * vy1 ) / (vx1 * vy2 - vy1 * vx2)
	// (sign of determinant:)
	// y = ((x10 - x20)*(vy1*vy2) - vx1 * y10 * vy2 + vx2 * y20 * vy1 ) / (vy1 * vx2 - vx1 * vy2)

	// x in [x_1, x_2]
	// y in [y_1, y_2]
	// D = vy1 * vx2 - vx1 * vy2
	//
	// x >= x_1 <=> (D > 0 && x(^) >= D*x_1 ) || (D < 0 && x(^) <= D*x_1)
	// ...
	// !!!
	// t > 0

	xUpper := (h2.p.y-h1.p.y)*h1.d.x*h2.d.x + h1.d.y*h1.p.x*h2.d.x - h2.d.y*h2.p.x*h1.d.x
	yUpper := (h1.p.x-h2.p.x)*h1.d.y*h2.d.y - h1.d.x*h1.p.y*h2.d.y + h2.d.x*h2.p.y*h1.d.y

	if D > 0 {
		if xUpper >= minCoord*D && xUpper <= maxCoord*D && yUpper >= minCoord*D && yUpper <= maxCoord*D {
			if h1.d.x > 0 {
				if xUpper < D*h1.p.x {
					return pos{}, fmt.Errorf("Time < 0 for %s", h1)
				}
			} else {
				if xUpper > D*h1.p.x {
					return pos{}, fmt.Errorf("Time < 0 for %s", h1)
				}
			}
			if h2.d.x > 0 {
				if xUpper < D*h2.p.x {
					return pos{}, fmt.Errorf("Time < 0 for %s", h2)
				}
			} else {
				if xUpper > D*h2.p.x {
					return pos{}, fmt.Errorf("Time <0 for %s", h2)
				}
			}
			return pos{x: xUpper / D, y: yUpper / D, z: 0}, nil
		}
		return pos{}, fmt.Errorf("Outside bbox D: %f , xUpper: %f yUpper: %f", D, xUpper, yUpper)
	}
	if D < 0 {
		if xUpper <= minCoord*D && xUpper >= maxCoord*D && yUpper <= minCoord*D && yUpper >= maxCoord*D {
			if h1.d.x > 0 {
				if xUpper > D*h1.p.x {
					return pos{}, fmt.Errorf("Time < 0 for %s", h1)
				}
			} else {
				if xUpper < D*h1.p.x {
					return pos{}, fmt.Errorf("Time < 0 for %s", h1)
				}
			}
			if h2.d.x > 0 {
				if xUpper > D*h2.p.x {
					return pos{}, fmt.Errorf("Time < 0 for %s", h2)
				}
			} else {
				if xUpper < D*h2.p.x {
					return pos{}, fmt.Errorf("Time <0 for %s", h2)
				}
			}
			return pos{x: xUpper / D, y: yUpper / D, z: 0}, nil
		}
		return pos{}, fmt.Errorf("Outside bbox D: %f , xUpper: %f yUpper: %f", D, xUpper, yUpper)
	}
	return pos{}, fmt.Errorf("Unreachable")
}

func (p1 pos) Minus(p2 pos) pos {
	return pos{x: p1.x - p2.x, y: p1.y - p2.y, z: p1.z - p2.z}
}

func (p1 pos) DotProduct(p2 pos) float64 {
	return p1.x*p2.x + p1.y*p2.y + p1.z*p2.z
}
func (p1 dir) Minus(p2 dir) dir {
	return dir{x: p1.x - p2.x, y: p1.y - p2.y, z: p1.z - p2.z}
}

func (p1 dir) DotProduct(p2 dir) float64 {
	return p1.x*p2.x + p1.y*p2.y + p1.z*p2.z
}

func (d1 dir) MultT(t float64) pos {
	return pos{x: d1.x * t, y: d1.y * t, z: d1.z * t}
}

func (p1 pos) DotProductDir(p2 dir) float64 {
	return p1.x*p2.x + p1.y*p2.y + p1.z*p2.z
}

func (p1 pos) Add(p2 pos) pos {
	return pos{x: p1.x + p2.x, y: p1.y + p2.y, z: p1.z + p2.z}
}
func (p1 dir) Add(p2 dir) dir {
	return dir{x: p1.x + p2.x, y: p1.y + p2.y, z: p1.z + p2.z}
}

// [(r1 + v1*t ) - ( r2 + v2*t )]^2 -> (min_t)
// 't
// [(r1 + v1*t) - (r2 + v2*t)].[v1 - v2] = 0
// .    [r2 - r1].[v1 - v2]
// t == -----------------
//           [v1 - v2]^2

func (h1 Hail) MinT(h2 Hail) float64 {
	if h1.d == h2.d {
		return 0.0
	}
	dd := h1.d.Minus(h2.d)
	return h2.p.Minus(h1.p).DotProductDir(dd) / dd.DotProduct(dd)
}

func (h1 Hail) MinDistance2(h2 Hail) float64 {
	t := h1.MinT(h2)
	if t < 0 {
		t = 0
	}
	p1 := h1.p.Add(h1.d.MultT(t))
	p2 := h2.p.Add(h2.d.MultT(t))
	dp := p1.Minus(p2)
	return dp.DotProduct(dp)
}

//      [r1 - rx].[vx - v1]
// t1 = -------------------
//          [vx - v1]^2
//
// r1 + v1*t1 = rx + vx*t1
// r2 + v2*t2 = rx + vx*t2
//
// (r1-rx) + (v1-vx) * t1 = 0
// (r2-rx) + (v2-vx) * t2 = 0
// r1 - rx == -t1 * (v1 - vx)
// r2 - rx == -t2 * (v2 - vx)
//
// (ax*bx + ay*by + az * bz) ^2 == (ax ^2 + ay^2 + az^2) * (bx^2 + by^2 + bz^2)

// vectors || <=> vector product is 0
// | i          j         k        |
// | (r1-rx)_x (r1-rx)_y (r1-rx)_z |
// | (v1-vx)_x (v1-vx)_y (v1-vx)_z |

// (r1-rx)_y * (v1-vx)_z - (r1-rx)_z * (v1-vx)_y = 0
// (r2-rx)_y * (v2-vx)_z - (r2-rx)_z * (v2-vx)_y = 0
// ...

// rxy vxz - rxz vxy -rxy v1z - r1y vxz + r1z vxy + rxz v1y + (r1y v1z - r1z v1y ) == 0
// rxy vxz - rxz vxy -rxy v2z - r2y vxz + r2z vxy + rxz v2y + (r2y v2z - r2z v2y ) == 0
// (-)
// -rxy (v2z - v1z) - vxz (r2z - r1z) + vxy (r2z - r1z) + rxz (v2y - v1y) + (r2y v2z - r2z v2y ) - (r1y v1z - r1z v1y ) == 0
// -rxy (v3z - v1z) - vxz (r3z - r1z) + vxy (r3z - r1z) + rxz (v3y - v1y) + (r3y v3z - r3z v3y ) - (r1y v1z - r1z v1y ) == 0

// (r1-rx)_z * (v1-vx)_x - (r1-rx)_x * (v1-vx)_z = 0
// (r2-rx)_z * (v2-vx)_x - (r2-rx)_x * (v2-vx)_z = 0
// ...
//
// rxz vxx - rxx vxz - rxz v1x - r1z vxx + r1x vxz + rxx v1z + (r1z v1x - r1x v1z) == 0
// rxz vxx - rxx vxz - rxz v2x - r2z vxx + r2x vxz + rxx v2z + (r2z v2x - r2x v2z) == 0
// (-)
// - rxz (v2x - v1x) - vxx (r2z - r1z) + vxz (r2x - r1x) + rxx (v2z - v1z) + (r2z v2x - r2x v2z) - (r1z v1x - r1x v1z) == 0
// - rxz (v3x - v1x) - vxx (r3z - r1z) + vxz (r3x - r1x) + rxx (v3z - v1z) + (r3z v3x - r3x v3z) - (r1z v1x - r1x v1z) == 0

// (r1-rx)_x * (v1-vx)_y - (r1-rx)_y * (v1-vx)_x = 0
// (r2-rx)_x * (v2-vx)_y - (r2-rx)_y * (v2-vx)_x = 0
// ...
//
// rxx vxy - rxy vxx - rxx v1y - r1x vxy + rxy v1x + r1y vxx + (r1x v1y - r1y v1x) == 0
// rxx vxy - rxy vxx - rxx v2y - r2x vxy + rxy v2x + r2y vxx + (r2x v2y - r2y v2x) == 0
// (-)
// - rxx (v2y - v1y) - vxy (r2x - r1x) + rxy (v2x - v1x) + vxx (r2y - r1y) + (r2x v2y - r2y v2x) - (r1x v1y - r1y v1x) == 0
// - rxx (v3y - v1y) - vxy (r3x - r1x) + rxy (v3x - v1x) + vxx (r3y - r1y) + (r3x v3y - r3y v3x) - (r1x v1y - r1y v1x) == 0

// System:

// - rxy (v2z - v1z) - vxz (r2z - r1z) + vxy (r2z - r1z) + rxz (v2y - v1y) + (r2y v2z - r2z v2y) - (r1y v1z - r1z v1y) == 0
// - rxy (v3z - v1z) - vxz (r3z - r1z) + vxy (r3z - r1z) + rxz (v3y - v1y) + (r3y v3z - r3z v3y) - (r1y v1z - r1z v1y) == 0
// - rxz (v2x - v1x) - vxx (r2z - r1z) + vxz (r2x - r1x) + rxx (v2z - v1z) + (r2z v2x - r2x v2z) - (r1z v1x - r1x v1z) == 0
// - rxz (v3x - v1x) - vxx (r3z - r1z) + vxz (r3x - r1x) + rxx (v3z - v1z) + (r3z v3x - r3x v3z) - (r1z v1x - r1x v1z) == 0
// - rxx (v2y - v1y) - vxy (r2x - r1x) + rxy (v2x - v1x) + vxx (r2y - r1y) + (r2x v2y - r2y v2x) - (r1x v1y - r1y v1x) == 0
// - rxx (v3y - v1y) - vxy (r3x - r1x) + rxy (v3x - v1x) + vxx (r3y - r1y) + (r3x v3y - r3y v3x) - (r1x v1y - r1y v1x) == 0

// rxx (v1y - v2y) + rxy (v2x - v1x) +                   vxx (r2y - r1y) + vxy (r1x - r2x)                     == (r1x v1y - r1y v1x) - (r2x v2y - r2y v2x)
// rxx (v1y - v3y) + rxy (v3x - v1x) +                   vxx (r3y - r1y) + vxy (r1x - r3x)                     == (r1x v1y - r1y v1x) - (r3x v3y - r3y v3x)
//                   rxy (v1z - v2z) + rxz (v2y - v1y) +                   vxy (r2z - r1z) + vxz (r1z - r2z)   == (r1y v1z - r1z v1y) - (r2y v2z - r2z v2y)
//                   rxy (v1z - v3z) + rxz (v3y - v1y) +                   vxy (r3z - r1z) + vxz (r1z - r3z)   == (r1y v1z - r1z v1y) - (r3y v3z - r3z v3y)
// rxx (v2z - v1z)                   + rxz (v1x - v2x) + vxx (r1z - r2z)                   + vxz (r2x - r1x)   == (r1z v1x - r1x v1z) - (r2z v2x - r2x v2z)
// rxx (v3z - v1z)                   + rxz (v1x - v3x) + vxx (r1z - r3z)                   + vxz (r3x - r1x)   == (r1z v1x - r1x v1z) - (r3z v3x - r3x v3z)

func (h1 Hail) At(t float64) pos {
	return pos{x: h1.p.x + h1.d.x*t, y: h1.p.y + h1.d.y*t, z: h1.p.z + h1.d.z*t}
}
func (inp Input) SolveAsSystem() Hail {
	v1 := inp.h[0].d
	r1 := inp.h[0].p
	v2 := inp.h[1].d
	r2 := inp.h[1].p
	v3 := inp.h[2].d
	r3 := inp.h[2].p

	a := mat.NewDense(6, 6, []float64{
		v1.y - v2.y, v2.x - v1.x, 0, r2.y - r1.y, r1.x - r2.x, 0,
		v1.y - v3.y, v3.x - v1.x, 0, r3.y - r1.y, r1.x - r3.x, 0,
		0, v1.z - v2.z, v2.y - v1.y, 0, r2.z - r1.z, r1.z - r2.z,
		0, v1.z - v3.z, v3.y - v1.y, 0, r3.z - r1.z, r1.z - r3.z,
		v2.z - v1.z, 0, v1.x - v2.x, r1.z - r2.z, 0, r2.x - r1.x,
		v3.z - v1.z, 0, v1.x - v3.x, r1.z - r3.z, 0, r3.x - r1.x,
	})

	b := mat.NewDense(6, 1, []float64{
		r1.x*v1.y - r1.y*v1.x - r2.x*v2.y + r2.y*v2.x,
		r1.x*v1.y - r1.y*v1.x - r3.x*v3.y + r3.y*v3.x,
		r1.y*v1.z - r1.z*v1.y - r2.y*v2.z + r2.z*v2.y,
		r1.y*v1.z - r1.z*v1.y - r3.y*v3.z + r3.z*v3.y,
		r1.z*v1.x - r1.x*v1.z - r2.z*v2.x + r2.x*v2.z,
		r1.z*v1.x - r1.x*v1.z - r3.z*v3.x + r3.x*v3.z,
	})

	var x mat.Dense
	var Ai = mat.DenseCopyOf(a)
	Ai.Inverse(a)
	err := x.Solve(a, b)
	if err != nil {
		log.Fatalf("no solution: %v", err)
	}
	var c mat.Dense
	c.Mul(Ai, b)
	fmt.Println(c)
	AR := mat.DenseCopyOf(a)
	AR.Mul(a, Ai)
	fmt.Println(AR)

	res := Hail{
		p: pos{
			x: c.At(0, 0),
			y: c.At(1, 0),
			z: c.At(2, 0),
		},
		d: dir{
			x: c.At(3, 0),
			y: c.At(4, 0),
			z: c.At(5, 0),
		},
	}
	t0 := res.MinT(inp.h[0])
	// (r1-rx)_x * (v1-vx)_y - (r1-rx)_y * (v1-vx)_x = 0
	v := (r1.x-res.p.x)*(v1.y-res.d.y) - (r1.y-res.p.y)*(v1.x-res.d.x)
	fmt.Printf("=== %f \n", v)
	P := mat.Dense{}

	fmt.Println(P)

	fmt.Printf("res: %s  ho: %s\n", res.At(t0), inp.h[0].At(t0))
	return res
}

// ...
//      [r2 - rx].[vx - v2]
// t2 = -------------------
//          [vx - v2]^2
//
// (r2 - rx + (v2 - vx) * t2)^2
//
// vx = (r1 + v1*t1 - rx)/t1 = (r1 - rx)/t1 + v1
//
// (r2 - rx + (v2 - v1 - (r1 - rx)/t1) * t2 ) ^2
// (r2 - rx(1-t2/t1) + (v2 - v1)t2 - r1(t2/t1) ) ^2  -> min_(rx_i, q, t2)
//
// t2/t1 = q
//
// (r2 - rx(1-q) + (v2 - v1)t2 - r1 q ) ^2
//
// (r2 - rx(1-q) + (v2 - v1)t2 - r1 q ).(rx - r1) == 0        // d/dq
// (r2 - rx(1-q) + (v2 - v1)t2 - r1 q ).(v2 - v1) == 0        // d/dt2

// t2 (v2-v1)^2 = (r1 q - r2 + rx*(1-q)).(v2-v1)
//
//      (r1 q - r2 + rx*(1-q)).(v2-v1)
// t2 = ------------------------------
//        (v2-v1)^2
//
// (r2 - rx(1-q) + (v2 - v1) (r1 q - r2 + rx*(1-q)).(v2-v1) / (v2-v1)^2 - r1 q ) . (rx - r1) == 0
//
// (r2 - rx + q (rx + (v2-v1) (r1 - rx).(v2-v1) /(v2-v1)^2 - r1)  + (v2-v1) (-r2 + rx).(v2-v1)/(v2-v1)^2 ) . (rx - r1) == 0

// q (rx + (v2-v1) (r1 - rx).(v2-v1) /(v2-v1)^2 - r1).(rx - r1) == (rx - r2 - (v2-v1) (-r2 + rx).(v2-v1)/(v2-v1)^2  ).(rx - r1)
//
// .    (rx - r2 - (v2-v1) (rx - r2).(v2-v1) / (v2-v1)^2  ).(rx - r1)
// q =  ------------------------------------------------------------
// .    (rx - r1 - (v2-v1) (rx - r1).(v2-v1) / (v2-v1)^2  ).(rx - r1)

//
//             (1-q)*(v2-v1)_x - rx_x * dq/drx_x * (v2-v1)_x
// dt2/drx_x = ---------------------------------------------
//             (v2-v1)^2
//
//
// d_q_upper/drx_x = (rx - r2 - (v2-v1) (rx - r2).(v2-v1) / (v2-v1)^2  )_x + rx_x *( 1 - (v2-v1)_x * (v2-v1)_x / (v2-v1)^2 )=
//                 = rx_x - r2_x - (v2-v1)_x (rx-r2).(v2-v1)/(v2-v1)^2 + rx_x - rx_x ( (v2-v1)_x^2/ (v2-v1)^2 )
//                 = rx_x (2 - (v2-v1)_x^2/(v2-v1)^2 ) - (v2-v1)_x (rx-r2).(v2-v1)/(v2-v1)^2 - r2_x
//                 = 2 rx_x (1 - (v2-v1)_x^2/(v2-v1)^2 ) - rx_y * (v2-v1)_x * (v2-v1)_y/(v2-v1)^2 -rx_z * (v2-v1)_x * (v2-v1)_z/(v2-v1)^2 + (v2-v1)_x r2.(v2-v1)/(v2-v1)^2 - r2_x
//
// d_q_lower/drx_x = rx_x (2 - (v2-v1)_x^2/(v2-v1)^2 ) - (v2-v1)_x (rx-r1).(v2-v1)/(v2-v1)^2 - r1_x
//
//            q_lower * d_q_upper/drx_x - q_upper * d_q_lower/drx_x
// dq/drx_x = -----------------------
//            (q_lower)^2
//
// (r2 - rx(1-q) + (v2 - v1)t2 - r1 q ) ^2    // d/drx_x == 0
//
// (r2 - rx(1-q) + (v2 - v1)t2 - r1 q ).{q-1 + (rx_x - r1_x) * dq/drx_x + (v2-v1)_x * dt2/drx_x , }

// r1 - rx = (vx - v1) * t1
// (r1 - rx) * [vx - v1]^2 = (vx - v1) * [r1 - rx].[vx - v1]
//
// (r1 - rx)_x * ( (vx - v1)_x^2 + (vx - v1)_y^2 + (vx-v1)_z^2 ) = (vx - v1)_x * ( (r1 - rx)_x*(vx - v1)_x + (r1-rx)_y * (vx - v1)_y + (r1 - rx)_z * (vx - v1)_z )
//                 ~~~~~~~~~~~~~                                                    ~~~~~~~~~~~~~~~~~~~~~~
// (r1-rx)_x *( (vx-v1)_y^2 + (vx-v1)_z^2 )  = (vx-v1)_x * ((r1-rx)_y * (vx-v1)_y + (r1-rx)_z * (vx-v1)_Z )
// R_x * (V_x^2 + V_y^2 + V_z^2) = V_x * (R_x * V_x + R_y * V_y + R_z * V_z)
// R_y * (V_x^2 + V_y^2 + V_z^2) = V_y * (R_x * V_x + R_y * V_y + R_z * V_z)
// R_z * (V_x^2 + V_y^2 + V_z^2) = V_z * (R_x * V_x + R_y * V_y + R_z * V_z)
//
// R_x * (V_y^2 + V_z^2) = V_x * (R_y * V_y + R_z * V_z)
// R_y * (V_x^2 + V_z^2) = V_y * (R_x * V_x + R_z * V_z)
// R_z * (V_x^2 + V_y^2) = V_z * (R_x * V_x + R_y * V_y)
//
// R_x V_y^2 + R_x V_z^2 - R_y V_x^2 - R_y V_z^2 - R_z V_x^2 - R_z V_y^2 =
// (R_x - R_z) V_y^2 + (R_x - R_y) V_z^2 - (R_y + R_z) V_x^2 =
//
//              = R_y V_x V_y + R_z V_x V_z - R_x V_x V_y - R_z V_y V_z - R_x V_x V_z - R_y V_y V_z
//              = (R_y - R_x) V_x V_y + (R_z - R_x) V_x V_z - (R_z + R_y) V_y V_z
//
// (R_x - R_z) (V_y^2 + V_x V_z) + (R_x - R_y) (V_z^2 + V_x V_y) - (R_y - R_z) (V_x^2 - V_y V_z) = 0

func (inp *Input) TotalFunctional(b Hail) float64 {
	s := float64(0)
	for _, h := range inp.h {
		h1 := Hail{
			p: pos{
				x: h.p.x / pfactor,
				y: h.p.y / pfactor,
				z: h.p.z / pfactor,
			},
			d: dir{
				x: h.d.x,
				y: h.d.y,
				z: h.d.z,
			},
		}
		s += h1.MinDistance2(b)
	}
	return math.Sqrt(s)
}
func (h1 Hail) Rev() Hail {
	return Hail{p: pos{x: -h1.p.x, y: -h1.p.y, z: -h1.p.z}, d: dir{x: -h1.d.x, y: -h1.d.y, z: -h1.d.z}}
}

func (h1 Hail) AddMult(axis Hail, x float64) Hail {
	return Hail{
		p: pos{
			x: h1.p.x + axis.p.x*x,
			y: h1.p.y + axis.p.y*x,
			z: h1.p.z + axis.p.z*x,
		},
		d: dir{
			x: h1.d.x + axis.d.x*x,
			y: h1.d.y + axis.d.y*x,
			z: h1.d.z + axis.d.z*x,
		},
	}
}
func (inp *Input) MinimiseAlongaxis(start Hail, axis Hail) (Hail, float64) {
	f := func(x float64) float64 {
		return inp.TotalFunctional(start.AddMult(axis, x))
	}
	tol := 1e-10
	maxIter := 500
	fnMaxFev := func(nfev int) bool { return nfev > 1500 }
	bm := optimize.NewBrentMinimizer(f, tol, maxIter, fnMaxFev)
	bm.Brack = []float64{-100, 100}
//	x, fx, nIter, nFev := bm.Optimize()
	x, fx, _, _ := bm.Optimize()
//	fmt.Printf("x: %.8g, fx: %.8g, nIter: %d, nFev: %d\n", x, fx, nIter, nFev)

	return start.AddMult(axis, x), fx
}
func (inp *Input) MinimiseAlongaxis1(start Hail, axis Hail) (Hail, float64) {
	h := start
	a := axis
	for {
		s0 := inp.TotalFunctional(h)
		if s0 == 0.0 {
			return h, s0
		}
		h1 := Hail{d: h.d.Add(a.d), p: h.p.Add(a.p)}
		s1 := inp.TotalFunctional(h1)
		if s1 == 0.0 {
			return h1, s0
		}
		axishalf := Hail{d: dir(a.d.MultT(0.5)), p: dir(a.p).MultT(0.5)}
		//		axistwice := Hail{d: dir(a.d.MultT(2.0)), p: dir(a.p).MultT(2.0)}
		h2 := Hail{d: h.d.Add(axishalf.d), p: h.p.Add(axishalf.p)}
		s2 := inp.TotalFunctional(h2)
		//
		// . s0   s2    s1
		//
		//
		switch {
		case math.Abs(s0-s1) < 1.0e-5 && math.Abs(s0-s2) < 1.0e-5:
			return h1, s1
		case s0 > s2 && s2 > s1:
			h = h1
			//			a = axistwice

		case s0 < s2 && s2 < s1:
			a = a.Rev()
			h = h2
		case s0 > s2 && s0 > s1 && s2 < s1:
			h = h1
			a = axishalf.Rev()
		case s0 <= s1 && s2 <= s0:
			h = h1
			a = axishalf

		}
	}
}

var pfactor = float64(1.0)

func (inp *Input) MinimizePowell() Hail {
	pm := optimize.NewPowellMinimizer()
	pm.Xtol = 0.000000001
	pm.Ftol = 0.0000001
	pm.MaxFev = 100000
	hopt := Hail{}
	pm.Callback = func(x []float64) {
		fmt.Printf("%.5f\n", x)
		hopt.p.x = x[0] * pfactor
		hopt.p.y = x[1] * pfactor
		hopt.p.z = x[2] * pfactor
		hopt.d.x = x[3]
		hopt.d.y = x[4]
		hopt.d.z = x[5]
	}
	pm.Logger = log.New(os.Stdout, "", 0)

	pm.Minimize(
		func(x []float64) float64 {
			h := Hail{p: pos{x: x[0], y: x[1], z: x[2]}, d: dir{x: x[3], y: x[4], z: x[5]}}
			s := inp.TotalFunctional(h)
			pm.Logger.Println(s)
			return s
		},
		[]float64{1, 1, 1, 1, 1, 1},
	)
	return hopt
}
func (inp *Input) MinimizeFull() Hail {

	res := Hail{}
	it := 1
	dx := float64(1.0)
	for {
		it += 1
		if it % 1000 == 0 {
			dx = dx/2.0
		}

		axis := Hail{d: dir{x: dx}}
		s := 0.0
		res, s = inp.MinimiseAlongaxis(res, axis)
		if math.Abs(s) < 1.0e-20 {
			return res
		}
		fmt.Printf("%d %f\n", it, s)
		axis = Hail{d: dir{y: dx}}
		res, s = inp.MinimiseAlongaxis(res, axis)
		if math.Abs(s) < 1.0e-20 {
			return res
		}
//		fmt.Printf("%f\n", s)
		axis = Hail{d: dir{z: dx}}
		res, s = inp.MinimiseAlongaxis(res, axis)
		if math.Abs(s) < 1.0e-20 {
			return res
		}
//		fmt.Printf("%f\n", s)

		axis = Hail{p: pos{x: dx}}
		res, s = inp.MinimiseAlongaxis(res, axis)
		if math.Abs(s) < 1.0e-20 {
			return res
		}
//		fmt.Printf("%f\n", s)
		axis = Hail{p: pos{y: dx}}
		res, s = inp.MinimiseAlongaxis(res, axis)
		if math.Abs(s) < 1.0e-20 {
			return res
		}
//		fmt.Printf("%f\n", s)
		axis = Hail{p: pos{z: dx}}
		res, s = inp.MinimiseAlongaxis(res, axis)
		if math.Abs(s) < 1.0e-20 {
			return res
		}
//		fmt.Printf("%f\n", s)
		fmt.Printf("%v\n", res)
	}
}

func (h1 Hail) GoesThroughBbox2D() bool {
	// y1(x1) = ( y10 * vx1 + vy1*(x1-x10) ) / vx1
	//
	//   o --- o
	//   |     |
	//   o --- o
	//
	y1upperMin := h1.p.y*h1.d.x + h1.d.y*(minCoord-h1.p.x)
	y1MinWhere := 0
	if h1.d.x > 0 {
		if y1upperMin > h1.d.x*maxCoord {
			// *
			// o -
			// |
			y1MinWhere = 1
		} else if y1upperMin < h1.d.x*minCoord {
			y1MinWhere = -1
		} else {
			return true
		}
	} else {
		if y1upperMin < h1.d.x*maxCoord {
			y1MinWhere = 1
		} else if y1upperMin > h1.d.x*minCoord {
			y1MinWhere = -1
		} else {
			return true
		}
	}
	y1upperMax := h1.p.y*h1.d.x + h1.d.y*(maxCoord-h1.p.x)
	if h1.d.x > 0 {
		if y1upperMax > h1.d.x*maxCoord {
			return y1MinWhere == -1
		} else if y1upperMax < h1.d.x*minCoord {
			return y1MinWhere == 1
		} else {
			return true
		}
	} else {
		if y1upperMax < h1.d.x*maxCoord {
			return y1MinWhere == -1
		} else if y1upperMax > h1.d.x*minCoord {
			return y1MinWhere == 1
		} else {
			return true
		}
	}
}

func (h1 Hail) String() string {
	return fmt.Sprintf("%f, %f , %f @ %f, %f, %f", h1.p.x, h1.p.y, h1.p.z , h1.d.x, h1.d.y, h1.d.z)
}

func (p pos) String() string {
	return fmt.Sprintf("[%f, %f, %f]", p.x, p.y, p.z)
}
func (inp Input) Count() (int, int) {
	s1 := 0
	s2 := 0
	goodh := []Hail{}
	for _, h := range inp.h {
		if h.GoesThroughBbox2D() {
			for _, h2 := range goodh {
				if p, err := h.Intersect2d(h2); err == nil {
					s1++
					fmt.Printf("%s %s with %s\n", p, h, h2)
				}
			}
			goodh = append(goodh, h)
		}
	}
	//ho := inp.MinimizePowell()
	ho := inp.MinimizeFull()
	//ho := inp.SolveAsSystem()
	s2 = int(ho.p.x + ho.p.y + ho.p.z)
	return s1, s2
}

func Readlines(file io.Reader) (int, int) {
	scanner := bufio.NewScanner(file)
	inp := ReadInput(scanner)
	return inp.Count()
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Ups")
	}
	defer file.Close()
	f, f1 := Readlines(file)
	fmt.Println(strconv.Itoa(f), f1)
}
