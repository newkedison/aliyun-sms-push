package main

func test_db() {
	//   ctx := getContextWithTimeout(5000)

	//   collection := dbClient.Database("sms_push").Collection("test1")
	//   res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
	//   if err != nil {
	//     panic(err)
	//   }
	//   id := res.InsertedID
	//   print("insertedID:", id)

	//   db := client.Database("sms_push")
	//   cur, err := db.ListCollections(getContextWithTimeout(1000), bson.D{})
	//   if err != nil {
	//     panic(err)
	//   }
	//   defer cur.Close(context.Background())
	//   for cur.Next(context.Background()) {
	//     // To decode into a struct, use cursor.Decode()
	//     result := struct {
	//       Name string
	//     }{}
	//     err := cur.Decode(&result)
	//     if err != nil {
	//       log.Fatal(err)
	//     }
	//     dump(result)
	//     // To get the raw bson bytes use cursor.Current
	//     raw := cur.Current
	//     // do something with raw...
	//     dump(raw)
	//   }
	//   if err := cur.Err(); err != nil {
	//     panic(err)
	//   }

	//   decimal, _ := primitive.ParseDecimal128("-1.2345")
	//   collection := dbClient.Database("sms_push").Collection("test_insert")
	//   count := 10000

	//   start := now()
	//   for i := 0; i < count; i++ {
	//     _, err := collection.InsertOne(context.Background(), data)
	//     if err != nil {
	//       panic(err)
	//     }
	//   }
	//   print(now().Sub(start).Seconds())
	//   time.Sleep(time.Millisecond * 100)
	//   n, err := collection.CountDocuments(context.Background(), &bson.D{})
	//   if err != nil {
	//     panic(err)
	//   }
	//   print(n)

	//   collection.Drop(getContextWithTimeout(1000))
	//   var arr []interface{}
	//   for i := 0; i < count; i++ {
	//     data := bson.M{
	//       "integer": i,
	//       "float":   float64(i) + 1.123,
	//       "decimal": decimal,
	//       "time":    primitive.DateTime(now().UnixNano() / 1000000),
	//     }
	//     arr = append(arr, data)
	//     time.Sleep(time.Microsecond * 10)
	//   }
	//   start := now()
	//   _, err = collection.InsertMany(getContextWithTimeout(10000), arr, options.InsertMany().SetOrdered(true))
	//   if err != nil {
	//     panic(err)
	//   }
	//   print(now().Sub(start).Seconds())
	//   time.Sleep(time.Millisecond * 200)
	//   n, err := collection.CountDocuments(context.Background(), &bson.D{})
	//   if err != nil {
	//     panic(err)
	//   }
	//   print(n)

	//   {
	//   type testStruct struct {
	//     ID     primitive.ObjectID `bson:"_id,omitempty"`
	//     A_A    int
	//     BB     float64
	//     CccCcc float32
	//     D      string
	//     E      time.Time
	//     F      primitive.Decimal128
	//   }
	//     collection := client.Database("sms_push").Collection("test_struct")
	//     collection.Drop(ctx)
	//     f, _ := primitive.ParseDecimal128("2.000000002")
	//     t := testStruct{
	//       A_A:    1,
	//       BB:     2.000000002,
	//       CccCcc: 2.000000002,
	//       D:      "abc",
	//       E:      now(),
	//       F:      f,
	//     }
	//     res, err := collection.InsertOne(ctx, t)
	//     if err != nil {
	//       panic(err)
	//     }
	//     dump(res)

	//     time.Sleep(100 * time.Millisecond)
	//     cursor, err := collection.Find(ctx, bson.D{})
	//     i := 0
	//     for cursor.Next(ctx) {
	//       i++
	//       var t2 testStruct
	//       err := cursor.Decode(&t2)
	//       if err != nil {
	//         panic(err)
	//       }
	//       dump(t2)
	//       print(t2.E.Local())
	//     }
	//     print(i)
	//   }

	//   {
	//     go func() {
	//       for {
	//         time.Sleep(3 * time.Second)
	//         dbClient.Disconnect(ctxEmpty)
	//         print("auto disconnect")
	//       }
	//     }()
	//     go func() {
	//       for {
	//         err := colSignName.FindOne(ctxEmpty, bson.D{}).Err()
	//         if err == mongo.ErrClientDisconnected {
	//           print("ErrClientDisconnected, reconnect")
	//           err = connectToDB(globalConfig.MongoDB.URI)
	//           if err != nil {
	//             panic(err)
	//           }
	//         } else {
	//           print("OK")
	//           time.Sleep(time.Second)
	//         }
	//       }
	//     }()
	//     select {}
	//   }

	//   {
	//     col := db.Collection("test_update")
	//     col.Drop(ctx)
	//     _, err := col.InsertOne(context.Background(), bson.M{
	//       "hello": "world",
	//       "aaa": bson.A{
	//         bson.D{
	//           {"aaa-1", 1},
	//           {"aaa-2", bson.A{1, 2, 3}},
	//         },
	//         bson.D{
	//           {"aaa-1", 2},
	//           {"aaa-2", bson.A{2, 3, 4}},
	//         },
	//       },
	//       "aaaa": bson.A{1, 2, 3},
	//     })
	//     if err != nil {
	//       panic(err)
	//     }
	//     res, err := col.UpdateOne(ctx, bson.D{{"hello", "world"}},
	//       bson.D{
	//         {"$currentDate", bson.D{
	//           {"lastModified", true},
	//         }},
	//       })
	//     dump(res)
	//     res, err = col.UpdateOne(ctx, bson.D{{"hello", "world"}},
	//       bson.D{
	//         {"$set", bson.D{{"bbb", "bbb"}}},
	//         {"$set", bson.D{{"aaa.$[element1].aaa-1", 99}}},
	//         {"$set", bson.D{{"aaa.$[element2].aaa-2.$[element3]", 99}}},
	//         {"$currentDate", bson.D{{"lastModified", true}}},
	//       },
	//       options.Update().SetArrayFilters(options.ArrayFilters{
	//         Registry: bson.DefaultRegistry,
	//         Filters: []interface{}{
	//           bson.D{{"element1.aaa-1", 2}},
	//           bson.D{{"element2.aaa-1", bson.D{{"$lt", 10}}}},
	//           bson.D{{"element3", bson.D{{"$gt", 1}}}},
	//         },
	//       }),
	//     )
	//     dump(res)
	//   }
	//   {
	//     ctx := getContextWithTimeout(5000)
	//     col := db.Collection("test_update")
	//     col.Drop(ctx)
	//     _, err := col.InsertOne(ctx, bson.M{"value": 1})
	//     if err != nil {
	//       panic(err)
	//     }
	//     _, err = col.InsertOne(ctx, bson.M{"value": 2})
	//     if err != nil {
	//       panic(err)
	//     }
	//     res, err := col.UpdateMany(ctx, bson.D{}, bson.D{
	//       {"$inc", bson.D{{"count", 1}}},
	//     })
	//     if err != nil {
	//       panic(err)
	//     }
	//     dump(res)
	//     res, err = col.UpdateMany(ctx, bson.D{{"value", 2}}, bson.D{
	//       {"$inc", bson.D{{"count", 1}}},
	//     })
	//     if err != nil {
	//       panic(err)
	//     }
	//     dump(res)
	//   }
}
