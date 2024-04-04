package treeprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGap(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddNode("hello")
	tree.AddGap()
	tree.AddNode("world")
	actual := tree.String()
	expected := `.
├── hello
│
└── world
`
	assert.Equal(expected, actual)
}

func TestGapWithChildren(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	n1 := tree.AddBranch("hello")
	n1.AddNode("foo")
	n1.AddNode("bar")
	tree.AddGap()
	n2 := tree.AddBranch("world")
	n2.AddNode("foo")
	n2.AddNode("bar")
	actual := tree.String()
	expected := `.
├── hello
│   ├── foo
│   └── bar
│
└── world
    ├── foo
    └── bar
`
	assert.Equal(expected, actual)
}

func TestNestedGap(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	n1 := tree.AddBranch("hello")
	n1.AddNode("foo")
	n1.AddGap()
	n1.AddNode("bar")
	actual := tree.String()
	expected := `.
└── hello
    ├── foo
    │
    └── bar
`
	assert.Equal(expected, actual)
}

func TestGapWithChildren2(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	n1 := tree.AddBranch("hello")
	tree.AddGap()
	n1.AddNode("foo")
	n1.AddNode("bar")
	tree.AddGap()
	n2 := tree.AddBranch("world")
	n2.AddNode("foo")
	n2.AddNode("bar")
	actual := tree.String()
	expected := `.
├── hello
│   ├── foo
│   └── bar
│
│
└── world
    ├── foo
    └── bar
`
	assert.Equal(expected, actual)
}
