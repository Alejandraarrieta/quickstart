package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(" "))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	/*err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)*/

	quickstartDatabase := client.Database("quickstart")
    podcastsCollection := quickstartDatabase.Collection("podcasts")
    episodesCollection := quickstartDatabase.Collection("episodes")

	//COMENTO LINEA 42 A 70 PARA PROBAR LA BUSQUEDA

	/*podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
		{Key: "tittle", Value:"Nadie dice nada"},
		{Key: "author", Value: "Nico"},
		{"tags", bson.A{"development","programing", "coding"}},
	}) 
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(podcastResult.InsertedID)

	episodesResult, err := episodesCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title","Episode #1"},
			{"description","This is the first episode"},
			{"duration",25},
		},
		bson.D{
			{"podcast", podcastResult.InsertedID},
			{"title","Episode #2"},
			{"description","This is the second episode"},
			{"duration",32},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(episodesResult.InsertedIDs)
	fmt.Printf("Inserted %v documents into episode collection!\n", len(episodesResult.InsertedIDs))*/

//Busqueda de todos los docs, 1ra forma.
	cursor, err	 := episodesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var episodes [] bson.M
	if err = cursor.All(ctx, &episodes); err != nil{
		log.Fatal(err)
	}
	fmt.Println(episodes)
}

// Busqueda por lote, es para muchos datos
cursor, err := episodesCollection.Find(ctx, bson.M{})
if err != nil {
    log.Fatal(err)
}
defer cursor.Close(ctx)
for cursor.Next(ctx) {
    var episode bson.M
    if err = cursor.Decode(&episode); err != nil {
        log.Fatal(err)
    }
    fmt.Println(episode)
}

//leer solo un doc
var podcast bson.M
if err = podcastsCollection.FindOne(ctx, bson.M{}).Decode(&podcast); err != nil {
    log.Fatal(err)
}
fmt.Println(podcast)

//Buscar con filtro
filterCursor, err := episodesCollection.Find(ctx, bson.M{"duration": 25})
if err != nil {
    log.Fatal(err)
}
var episodesFiltered []bson.M
if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
    log.Fatal(err)
}
fmt.Println(episodesFiltered)


//Buscar estableciendo orden con config de opciones
opts := options.Find()
opts.SetSort(bson.D{{"duration", -1}})
sortCursor, err := episodesCollection.Find(ctx, bson.D{{"duration", bson.D{{"$gt", 24}}}}, opts)
if err != nil {
    log.Fatal(err)
}
var episodesSorted []bson.M
if err = sortCursor.All(ctx, &episodesSorted); err != nil {
    log.Fatal(err)
}
fmt.Println(episodesSorted)

