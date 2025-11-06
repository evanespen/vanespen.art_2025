use bytes::Bytes;
use futures_util::StreamExt;

#[derive(Debug)]
struct NatsMessage {
    key: String,
}

#[tokio::main]
async fn main() -> Result<(), async_nats::Error> {
    // Connect to the NATS server
    let client = async_nats::connect("localhost:4222").await?;
    println!("NATS connected");

    // Subscribe to the "messages" subject
    let mut subscriber = client.subscribe("pictures").await?;
    println!("Subscribed");

    // Receive and process messages
    while let Some(message) = subscriber.next().await {

	let toto: NatsMessage = rmp_serde::from_slice(&message.payload).unwrap();
	println!("{:?}", toto);
	
        println!("Received message {:?}", message.payload);
    }

    Ok(())
}
