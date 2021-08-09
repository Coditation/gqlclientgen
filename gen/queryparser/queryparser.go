package queryparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
)

func ParseQueryDocuments(queryPath string, parseGql *ast.Schema) (*ast.QueryDocument, error) {
	queryDir, err := filepath.Abs(queryPath)
	if err != nil {
		return nil, err
	}
	dir, err := os.ReadDir(queryDir)
	if err != nil {
		return nil, err
	}
	querySources := make([]*ast.Source, 0, len(dir))
	for _, fileName := range dir {
		if filepath.Ext(fileName.Name()) == ".graphql" {
			var err error
			var b []byte
			b, err = ioutil.ReadFile(path.Join(queryDir, fileName.Name()))
			if err != nil {
				return nil, err
			}
			querySources = append(querySources, &ast.Source{Name: fileName.Name(), Input: string(b)})
		}
	}
	merger := newMerger()
	for _, source := range querySources {
		querySchema, err := parser.ParseQuery(source)
		if err != nil {
			return nil, err
		}
		merger.mergeQueryDocument(querySchema)
	}
	if errs := validator.Validate(parseGql, &merger.document); errs != nil {
		return nil, fmt.Errorf(": %w", errs)
	}

	return &merger.document, err
}

type merger struct {
	document      ast.QueryDocument
	unamedIndex   int
	unamedPattern string
}

func newMerger() *merger {
	unamedPattern := "Unamed"
	return &merger{unamedPattern: unamedPattern}
}

func (m *merger) mergeQueryDocument(other *ast.QueryDocument) {
	for _, operation := range other.Operations {
		if operation.Name == "" {
			m.unamedIndex++
			operation.Name = fmt.Sprintf("%s%d", m.unamedPattern, m.unamedIndex)
		}
	}
	m.document.Operations = append(m.document.Operations, other.Operations...)
	m.document.Fragments = append(m.document.Fragments, other.Fragments...)
}

func QueryDocumentsByOperations(schema *ast.Schema, operations ast.OperationList) ([]*ast.QueryDocument, error) {
	queryDocuments := make([]*ast.QueryDocument, 0, len(operations))
	for _, operation := range operations {
		fragments := fragmentsInOperationDefinition(operation)

		queryDocument := &ast.QueryDocument{
			Operations: ast.OperationList{operation},
			Fragments:  fragments,
			Position:   nil,
		}

		if errs := validator.Validate(schema, queryDocument); errs != nil {
			return nil, fmt.Errorf(": %w", errs)
		}

		queryDocuments = append(queryDocuments, queryDocument)
	}

	return queryDocuments, nil
}

func fragmentsInOperationDefinition(operation *ast.OperationDefinition) ast.FragmentDefinitionList {
	fragments := fragmentsInOperationWalker(operation.SelectionSet)
	uniqueFragments := fragmentsUnique(fragments)

	return uniqueFragments
}

func fragmentsUnique(fragments ast.FragmentDefinitionList) ast.FragmentDefinitionList {
	uniqueMap := make(map[string]*ast.FragmentDefinition)
	for _, fragment := range fragments {
		uniqueMap[fragment.Name] = fragment
	}

	uniqueFragments := make(ast.FragmentDefinitionList, 0, len(uniqueMap))
	for _, fragment := range uniqueMap {
		uniqueFragments = append(uniqueFragments, fragment)
	}

	return uniqueFragments
}

func fragmentsInOperationWalker(selectionSet ast.SelectionSet) ast.FragmentDefinitionList {
	var fragments ast.FragmentDefinitionList
	for _, selection := range selectionSet {
		var selectionSet ast.SelectionSet
		switch selection := selection.(type) {
		case *ast.Field:
			selectionSet = selection.SelectionSet
		case *ast.InlineFragment:
			selectionSet = selection.SelectionSet
		case *ast.FragmentSpread:
			fragments = append(fragments, selection.Definition)
			selectionSet = selection.Definition.SelectionSet
		}

		fragments = append(fragments, fragmentsInOperationWalker(selectionSet)...)
	}

	return fragments
}
