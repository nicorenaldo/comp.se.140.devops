import http.server
import json
import logging
import os
import signal
import socketserver
import subprocess
import sys

PORT = int(os.getenv("PORT", 8199))


def get_ip() -> str:
    return subprocess.check_output(["hostname", "-I"]).decode("utf-8").strip()


def get_ps() -> list[str]:
    return subprocess.check_output(["ps", "aux"]).decode("utf-8").strip().split("\n")


def get_df() -> str:
    return subprocess.check_output(["df", "-h"]).decode("utf-8").split("\n")


def get_boot() -> str:
    return " ".join(
        subprocess.check_output(["last", "reboot", "|", "tail", "-1"])
        .decode("utf-8")
        .strip()
        .split(" ")[2:]
    )


class GetInfoHandler(http.server.SimpleHTTPRequestHandler):
    def get_info(self) -> dict:
        ip = get_ip()
        ps = get_ps()
        df = get_df()
        last_boot = get_boot()

        return {
            "ip": ip,
            "ps": ps,
            "df": df,
            "last_boot": last_boot,
        }

    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/plain")
        self.end_headers()

        payload = self.get_info()
        self.wfile.write(json.dumps(payload, indent=4).encode("utf-8"))


Handler = GetInfoHandler


def signal_handler(sig, frame):
    logging.info("\nShutting down the server gracefully...")
    sys.exit(0)


if __name__ == "__main__":
    Handler = GetInfoHandler
    with socketserver.TCPServer(("", PORT), Handler) as httpd:
        signal.signal(signal.SIGINT, signal_handler)  # Handle Ctrl+C
        logging.info(f"Serving at port {PORT}")
        httpd.serve_forever()
