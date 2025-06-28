import { MongoClient } from "mongodb";
import * as dotenv from "dotenv";
dotenv.config();

const { MONGODB_URI, MONGODB_DB = "eventsdb" } = process.env;

const client = new MongoClient(MONGODB_URI);
export const dbReady = client.connect().then(() => {
  console.log("🟢  Mongo connected");
  return client.db(MONGODB_DB).collection("events");
});
