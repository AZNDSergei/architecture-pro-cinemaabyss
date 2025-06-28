import { Kafka } from "kafkajs";
import * as dotenv from "dotenv";
dotenv.config();

const {
  CLIENT_ID,
  GROUP_ID,
  KAFKA_BROKERS,
  KAFKA_TOPIC_MOVIE,
  KAFKA_TOPIC_USER,
  KAFKA_TOPIC_PAYMENT
} = process.env;

const kafka = new Kafka({
  clientId: CLIENT_ID,
  brokers: KAFKA_BROKERS.split(",")
});

export const producer = kafka.producer();
export const consumer = kafka.consumer({ groupId: GROUP_ID });

export const topics = {
  movie: KAFKA_TOPIC_MOVIE,
  user: KAFKA_TOPIC_USER,
  payment: KAFKA_TOPIC_PAYMENT
};

export async function initKafka(onMessage) {
  await producer.connect();
  await consumer.connect();

  await Promise.all(
    Object.values(topics).map((t) => consumer.subscribe({ topic: t, fromBeginning: true }))
  );

  consumer.run({
    eachMessage: async ({ topic, message }) => onMessage(topic, message)
  });

  console.log("Kafka connected & consumer started");
}
