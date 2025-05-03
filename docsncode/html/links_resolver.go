package html

import (
	"log"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type linksResolverTransformer struct{}

func traverseChildren(node ast.Node) {
	if node == nil {
		return
	}

	if node.Kind() == ast.KindLink {
		link := node.(*ast.Link)
		log.Printf("Found link, dest=%s", link.Destination)
		// TODO: update link destination

	}

	if node.HasChildren() {
		traverseChildren(node.FirstChild())
	}
	traverseChildren(node.NextSibling())
}

func (*linksResolverTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	traverseChildren(node)
}

var LinksResolverTransformer = &linksResolverTransformer{}
