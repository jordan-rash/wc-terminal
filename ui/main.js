import "./style.css";
import Alpine from "alpinejs";
import {
  connect,
  consumerOpts,
  headers,
  ErrorCode,
  credsAuthenticator,
  JSONCodec,
  StringCodec,
} from "nats.ws";

Alpine.data("terminal", (subject) => ({
  id: Math.random().toString(36).slice(2, 10),
  data: [],
  command_data: [],
  chat_data: [],
  mode: "command",
  currentnode: "0",
  nats: null,
  jc: null,
  sc: null,
  h12: false,
  modes: new Map([
    [
      "command",
      {
        name: "command",
        subject_postfix: "command",
        icon: {
          v: "0 0 256 256",
          d:
            "M119.97217,136.96875l-72,64a11.99975,11.99975,0,1,1-15.94434-17.9375L93.9375,128,32.02783,72.96875a11.99975,11.99975,0,1,1,15.94434-17.9375l72,64a11.99925,11.99925,0,0,1,0,17.9375ZM215.99414,180h-96a12,12,0,0,0,0,24h96a12,12,0,1,0,0-24Z",
        },
        placeholder: "Enter command...? for help",
      },
    ],
    [
      "chat",
      {
        name: "chat",
        subject_postfix: "chat",
        icon: {
          v: "0 0 512 512",
          d:
            "M385.045,0H126.955C56.96,0,0,56.96,0,126.976v130.048C0,327.04,56.96,384,126.955,384h22.379v106.667,c0,10.048,6.997,18.709,16.811,20.843c1.515,0.341,3.029,0.491,4.523,0.491c8.235,0,15.893-4.757,19.413-12.48,C218.645,436.8,286.4,394.347,304.235,384h80.811C455.04,384,512,327.04,512,257.024V126.976C512,56.96,455.04,0,385.045,0z",
        },
        placeholder: "Enter chat text...? for help",
      },
    ],
  ]),
  icons: new Map([
    [
      "directory",
      {
        icon: {
          v: "0 0 512 512",
          d:
            "M490.053,118.398H259.926V77.852c0-12.121-9.826-21.947-21.947-21.947H71.731c-12.121,0-21.947,9.826-21.947,21.947,v40.546H21.947C9.826,118.398,0,128.224,0,140.345v293.804c0,12.121,9.826,21.947,21.947,21.947h468.106,c12.121,0,21.947-9.826,21.947-21.947V140.345C512,128.224,502.174,118.398,490.053,118.398z M93.678,99.799h122.354v18.599,H93.678V99.799z M468.106,412.201H43.894v-249.91h424.212V412.201z",
        },
      },
    ],
    [
      "file",
      {
        icon: {
          v: "0 0 293.151 293.151",
          d:
            "M255.316,55.996l-51.928-52.94C201.471,1.102,198.842,0,196.104,0h-82.302h-8.232H45.113,c-5.631,0-10.197,4.566-10.197,10.192c0,5.626,4.566,10.192,10.197,10.192h60.457h8.232h72.11l0.005,44.231,c0,5.631,4.561,10.197,10.192,10.197h41.731v197.955H56.592V47.828c0-5.631-4.566-10.197-10.197-10.197,c-5.631,0-10.192,4.566-10.192,10.197v235.131c0,5.631,4.566,10.192,10.192,10.192h201.642c5.631,0,10.197-4.566,10.197-10.192,V63.137C258.229,60.467,257.185,57.903,255.316,55.996z M206.307,54.423V35.147l18.906,19.276H206.307z",
        },
      },
    ],
  ]),
  async init() {
    const server = import.meta.env.VITE_NATS_SERVER;
    //const nats_creds = import.meta.env.VITE_NATS_CREDS;
    const nats_creds = `-----BEGIN NATS USER JWT-----
    <USER JWT HERE>
------END NATS USER JWT------

************************* IMPORTANT *************************
NKEY Seed printed below can be used to sign and prove identity.
NKEYs are sensitive and should be treated as secrets.

-----BEGIN USER NKEY SEED-----
<USER SEED HERE>
------END USER NKEY SEED------

*************************************************************
`;

    this.jc = JSONCodec();
    this.sc = StringCodec();
    this.data = this.command_data;

    try {
      this.nats = await connect({
        servers: server,
        authenticator: credsAuthenticator(new TextEncoder().encode(nats_creds)),
      });
    } catch (e) {
      console.log("Error occurred", e);
    }

    const opts = consumerOpts();
    opts.orderedConsumer();
    const sub = await this.nats.jetstream().subscribe(subject, opts);

    try {
      for await (const m of sub) {
        const data = this.jc.decode(m.data);
        console.log(data);
        switch (data.type) {
          case "command":
            if (data.id !== this.id) {
              this.command_data.push(data);
            }
            break;
          case "chat":
            if (data.id !== this.id) {
              this.chat_data.push(data);
            }
            break;
          case "clear":
            if (this.mode == "command") {
              this.command_data.length = 0;
            }
            if (this.mode == "chat") {
              this.chat_data.length = 0;
            }
            this.data.length = 0;
            break;
          default:
            break;
        }
      }
    } catch (e) {
      console.log("Error occurred", e);
    }
  },

  async terminalInput(e) {
    let userCmd = e.view.document.getElementById("userinput").value.split(" ");
    let slash = userCmd[0].charAt(0);
    if (slash != "/") {
      switch (userCmd[0]) {
        case "":
          break;
        case "clear":
          this.data.length = 0;
          this.clear();
          break;
        case "?":
          // available slash commands
          // available terminal commands
          // available modes
          break;
        default:
          const msg = {
            id: this.id,
            ts: new Date(),
            type: this.mode,
            action: userCmd[0],
            fsname: "fs1",
            session: subject,
            payload: { nodeid: this.currentnode, eventNodeKey: userCmd[1] },
            resp: "",
          };
          switch (this.mode) {
            case "command":
              console.log(msg);
              try {
                const m = await this.nats.request(
                  subject + "." + this.modes.get(this.mode).subject_postfix,
                  this.jc.encode(msg)
                );
                console.dir(this.jc.decode(m.data));
                msg.resp = this.jc.decode(m.data);
                this.currentnode = msg.resp.currnodeid;
              } catch (err) {
                switch (err.code) {
                  case ErrorCode.NoResponders:
                    console.log("no one is listening to 'hello.world'");
                    break;
                  case ErrorCode.Timeout:
                    console.log("someone is listening but didn't respond");
                    break;
                  default:
                    console.log("request failed", err);
                }
              }
              break;
            case "chat":
              msg.action = e.view.document.getElementById("userinput").value;
              break;
          }
          this.nats.publish(subject, this.jc.encode(msg));
          this.data.push(msg);
      }
    } else {
      const msg = {
        id: this.id,
        ts: new Date(),
        type: "slash",
        req: userCmd,
        resp: "",
      };
      let t = e.view.document.getElementById("termicon");
      let i = e.view.document.getElementById("termpath");
      let p = e.view.document.getElementById("userinput");
      switch (userCmd[0]) {
        case "/nick":
          if (userCmd[1] === "") {
            break;
          } else {
            this.id = userCmd[1];
          }
          break;
        case "/cmd":
          if (this.mode != "command") {
            this.mode = "command";
            this.data = this.command_data;
            msg.resp = "Entering command mode";
          }
          break;
        case "/chat":
          if (this.mode != "chat") {
            this.mode = "chat";
            this.data = this.chat_data;
            msg.resp = "Entering chat mode";
          }
          break;
        default:
          msg.resp = "Slash command not supported";
      }
      t.setAttribute("viewBox", this.modes.get(this.mode).icon.v);
      i.setAttribute("d", this.modes.get(this.mode).icon.d);
      p.setAttribute("placeholder", this.modes.get(this.mode).placeholder);
      p.setAttribute(
        "onblur",
        'this.placeholder="' + this.modes.get(this.mode).placeholder + '"'
      );

      if (msg.resp != "") {
        this.data.push(msg);
      }
    }

    e.view.document.getElementById("userinput").value = "";
  },

  clear() {
    const msg = { id: this.id, type: "clear" };
    const h = headers();
    h.set("Nats-Rollup", "sub");
    this.nats.publish(subject, this.jc.encode(msg), { headers: h });
  },

  getDisplayTime(t) {
    const options = {
      hour12: this.h12,
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    };
    let current = new Date(t);
    return "[" + current.toLocaleString("en-US", options) + "]";
  },

  getDisplay(c) {
    return "> " + c;
  },

  getPlaceholder() {
    switch (this.mode) {
      case "command":
        return "Enter command....? for help";
      case "chat":
        return "Enter chat message....? for help";
    }
  },

  getIcon(y) {
    console.log(this.icons.get(y));
    return this.icons.get(y);
  },

  getBool(b) {
    if (b === "true") {
      return true;
    } else {
      return false;
    }
  },
  guidGenerator() {
    var S4 = function () {
      return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
    };
    return (
      S4() +
      S4() +
      "-" +
      S4() +
      "-" +
      S4() +
      "-" +
      S4() +
      "-" +
      S4() +
      S4() +
      S4()
    );
  },
}));

Alpine.start();
