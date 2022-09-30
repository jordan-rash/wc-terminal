use rust_embed::RustEmbed;
use wasmbus_rpc::actor::prelude::*;
use wasmcloud_interface_httpserver::{
    HeaderMap, HttpRequest, HttpResponse, HttpServer, HttpServerReceiver,
};
use wasmcloud_interface_logging::debug;

#[derive(RustEmbed)]
#[folder = "./dist"]
struct Asset;

#[derive(Debug, Default, Actor, HealthResponder)]
#[services(Actor, HttpServer)]
struct UiActor {}

#[async_trait]
impl HttpServer for UiActor {
    async fn handle_request(
        &self,
        _ctx: &Context,
        req: &HttpRequest,
    ) -> std::result::Result<HttpResponse, RpcError> {
        let path = req.path.to_string();
        let trimmed = if path.trim() == "/" {
            // Default to index.html if the root path is given alone
            debug!("Found root path, assuming index.html");
            "index.html"
        } else {
            path.trim().trim_start_matches('/')
        };

        debug!("Got path {}, attempting to fetch", trimmed);
        if let Some(file) = Asset::get(trimmed) {
            debug!(
                "Found file {}, returning {} bytes",
                trimmed,
                file.data.len()
            );

            let content_type = mime_guess::from_path(trimmed)
                .first()
                .map(|m| m.to_string());

            let mut header = HeaderMap::new();
            if let Some(content_type) = content_type {
                debug!(
                    "Found content type of {}, setting Content-Type header",
                    content_type
                );
                header.insert("Content-Type".to_string(), vec![content_type]);
            }

            return Ok(HttpResponse {
                body: Vec::from(file.data),
                header,
                ..Default::default()
            });
        };

        debug!("Did not find file {}, returning", trimmed);
        return Ok(HttpResponse {
            ..Default::default()
        });
    }
}
