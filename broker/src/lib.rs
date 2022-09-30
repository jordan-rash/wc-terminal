use wasmbus_rpc::actor::prelude::*;
use wasmcloud_interface_logging::{debug, info};
use wasmcloud_interface_messaging::{
    MessageSubscriber, MessageSubscriberReceiver, Messaging, MessagingSender, PubMessage,
    SubMessage,
};

mod naryfs;
use naryfs::*;

const FS_ACTOR: &str = "MDW6PAE7PLGJAWUMWXGFF4TUBCBQ3NHUMVRZA53EBTHRFH3HZV32XICZ";
#[derive(Debug, Default, Actor, HealthResponder)]
#[services(Actor, MessageSubscriber)]
struct BrokerActor {}

/// Implementation of HttpServer trait methods
#[async_trait]
impl MessageSubscriber for BrokerActor {
    async fn handle_message(&self, ctx: &Context, msg: &SubMessage) -> RpcResult<()> {
        debug!("{:?}", msg);
        let action = serde_json::from_slice::<naryfs::FsMsg>(&msg.body)
            .map_err(|e| tag_err("decoding metadata", e))?;

        let sender = naryfs::FsSubscriberSender::to_actor(FS_ACTOR);
        let resp = sender.handle_fs_message(ctx, &action).await?;
        let serialized_resp = serde_json::to_string(&resp).unwrap();

        let provider = MessagingSender::new();
        let _ = provider
            .publish(
                ctx,
                &PubMessage {
                    body: serialized_resp.as_bytes().to_vec(),
                    reply_to: None,
                    subject: msg.reply_to.as_ref().unwrap().to_string(),
                },
            )
            .await;
        info!("BROKER RESPONSE: {:?}", serialized_resp.to_string());

        Ok(())
    }
}

fn tag_err<T: std::string::ToString>(msg: &str, e: T) -> RpcError {
    RpcError::ActorHandler(format!("{}: {}", msg, e.to_string()))
}
