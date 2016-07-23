// Copyright 2016 The Gini Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package logic_test

import (
	"gini/logic"
	"gini/z"
	"log"
	"math/rand"
	"testing"
)

func TestCGrowStrash(t *testing.T) {
	c := logic.NewC()
	N := 1020
	ins := make([]z.Lit, 0, N)
	for i := 0; i < N; i++ {
		ins = append(ins, c.NewIn())
	}
	gs := make([]z.Lit, N/2)
	for i := 0; i < N/2; i++ {
		j := len(ins) - 1 - i
		a, b := ins[i], ins[j]
		g := c.And(a, b)
		gs[i] = g
	}
	for i := 0; i < N/2; i++ {
		j := len(ins) - 1 - i
		a, b := ins[i], ins[j]
		g := c.And(a, b)
		if g != gs[i] {
			t.Errorf("invalid strash")
		}
	}
}

type op struct {
	a z.Lit
	b z.Lit
	g z.Lit
}

func TestCLogic(t *testing.T) {
	c := logic.NewC()
	a := c.NewIn()
	b := c.NewIn()
	ops := []op{
		op{a: c.T, b: c.NewIn()},
		op{a: c.F, b: c.NewIn()},
		op{a: a, b: a},
		op{a: a, b: a.Not()},
		op{a: a, b: b},
		op{a: b, b: a},
		op{a: c.NewIn(), b: c.NewIn()}}

	for i, _ := range ops {
		ops[i].g = c.And(ops[i].a, ops[i].b)
	}
	if ops[0].g != ops[0].b {
		t.Errorf("t simp")
	}
	if ops[1].g != ops[1].a {
		t.Errorf("f simp")
	}
	if ops[2].g != ops[2].a {
		t.Errorf("= simp")
	}
	if ops[3].g != c.F {
		t.Errorf("!= simp")
	}
	if ops[4].g != ops[5].g {
		t.Errorf("h simp")
	}
}

func TestEval(t *testing.T) {
	c := logic.NewC()
	a, b := c.NewIn(), c.NewIn()
	g := c.And(a, b)
	_ = g
	log.Printf("%s = %s and %s\n", g, a, b)
	vs := make([]bool, 5)
	vs[2], vs[3] = true, true
	log.Printf("b4 %+v\n", vs)
	c.Eval(vs)
	log.Printf("after %+v\n", vs)

	if vs[4] != true {
		t.Errorf("bad and eval")
	}
}

var rnd = rand.New(rand.NewSource(1))

func TestEval64(t *testing.T) {
	c := logic.NewC()
	a, b := c.NewIn(), c.NewIn()
	c.And(a, b)
	vs := make([]uint64, 5)
	for i := 0; i < 5; i++ {
		vs[i] = uint64(rnd.Int63())
	}
	c.Eval64(vs)
	for i := 0; i < 63; i++ {
		s := uint64(1 << uint64(i))
		a := (vs[2] & s) != 0
		b := (vs[3] & s) != 0
		c := (vs[4] & s) != 0
		if a && b && !c {
			t.Errorf("not true")
		} else if (!a || !b) && c {
			t.Errorf("not false")
		}
	}
}
