package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
)

type Tutorial struct {
	ID          int
	Titulo      string
	Autor       Autor
	Comentarios []Comentario
}

type Autor struct {
	Nombre     string
	Tutoriales []int
}

type Comentario struct {
	Cuerpo string
}

func popular() []Tutorial {
	autor := &Autor{Nombre: "juan sains", Tutoriales: []int{1}}
	tutorial := Tutorial{
		ID:     1,
		Titulo: "Tutorial basico Go y Graphql",
		Autor:  *autor,
		Comentarios: []Comentario{
			Comentario{Cuerpo: "Primer comentario"},
		},
	}
	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)

	return tutorials
}

func main() {
	fmt.Println("tutorial basico con graphql")
	tutoriales := popular()

	var tipoComentario = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comentario",
			Fields: graphql.Fields{
				"cuerpo": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var tipoAutor = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Autor",
			Fields: graphql.Fields{
				"Nombre": &graphql.Field{
					Type: graphql.String,
				},
				"Tutoriales": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)

	var tipoTutorial = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Tutorial",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"titulo": &graphql.Field{
					Type: graphql.String,
				},
				"autor": &graphql.Field{
					Type: tipoAutor,
				},
				"comentarios": &graphql.Field{
					Type: graphql.NewList(tipoComentario),
				},
			},
		},
	)

	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tipoTutorial,
			Description: "obtener el tutorial por el ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, tutorial := range tutoriales {
						if int(tutorial.ID) == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"lista": &graphql.Field{
			Type:        graphql.NewList(tipoTutorial),
			Description: "obtener toda la lista de tutoriales",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutoriales, nil
			},
		},
	}

	/* ejemplo 1
	fields := graphql.Fields{
		"hola": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error){
				return"mundo", nil
			},
		},
	}*/
	//define la configuracion del objeto
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	//define la configuracion del esquema
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	//creamos nuestro esquema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("fallo al crear el nuevo esquema GraphQL, err %v", err)
	}
	query := `
		{
			tutorial(id:1){
				titulo
				autor{
					Nombre
					Tutoriales
				}
			}
		}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("fallo al ejecutar la operacion graphql, error: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
