import express from "express";
import * as dotenv from "dotenv";
import { initKafka, producer, topics } from "./kafka.js";
import { dbReady } from "./db.js";
import { logEvent } from "./logger.js";

dotenv.config();
const app = express();
const PORT = process.env.PORT || 8082;

let ready = false;

app.use(express.json());

app.get("/api/events/health", (_, res) =>
  res.status(ready ? 200 : 503).json({ status: ready })
);

function eventHandler(topicKey) {
  return async (req, res, next) => {
    try {
      const payload = { type: topicKey, ts: new Date(), data: req.body };
      await producer.send({
        topic: topics[topicKey],
        messages: [{ value: JSON.stringify(payload) }]
      });
      logEvent("out", topics[topicKey], payload);
      res.status(201).json({ status: "success" });
    } catch (err) {
      next(err);
    }
  };
}

app.post("/api/events/movie", eventHandler("movie"));
app.post("/api/events/user", eventHandler("user"));
app.post("/api/events/payment", eventHandler("payment"));

app.use((err, _req, res, _next) =>
  res.status(500).json({ error: err.message })
);

// 👉 запуск сервера только после готовности
(async () => {
  try {
    const collection = await dbReady;

    await initKafka(async (topic, { value }) => {
      const doc = JSON.parse(value.toString());
      logEvent("in", topic, doc);
      await collection.insertOne({ topic, ...doc });
    });

    ready = true;
    console.log("Event-service ready");

    app.listen(PORT, () =>
      console.log(`✅ API listening on port ${PORT}`)
    );
  } catch (err) {
    console.error("❌ Startup error:", err);
    process.exit(1);
  }
})();
