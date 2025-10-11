#!/usr/bin/env bash
set -euo pipefail

API="${API:-http://localhost:8080}"
PASS="${PASS:-secret123}"  # usar algo >= 8 chars

say() { printf "\n=== %s ===\n" "$*"; }

login() {
  local email="$1" pass="$2"
  curl -s -o /dev/stderr -w "\n[HTTP %{http_code}]\n" \
    -X POST "$API/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"email\":\"$email\",\"password\":\"$pass\"}"
}

register() {
  local name="$1" email="$2" pass="$3"
  say "login $email (debería fallar si no existe)"
  login "$email" "$pass" || true

  say "registrando $email"
  curl -s -o /dev/stderr -w "\n[HTTP %{http_code}]\n" \
    -X POST "$API/auth/register" \
    -H 'Content-Type: application/json' \
    -d "{\"name\":\"$name\",\"email\":\"$email\",\"password\":\"$pass\"}"
}

say "healthz"
curl -i "$API/healthz" || true

# Usuarios (todos con la misma clave)
register "Maestro Roshi"         "roshi@kamehouse.example"         "$PASS"
register "Krillin"               "krillin@kamehouse.example"       "$PASS"
register "Yamcha"                "yamcha@capsule.example"          "$PASS"
register "Goku"                  "goku@capsule.example"            "$PASS"

register "Baki Hanma"            "baki.hanma@example.example"              "$PASS"
register "Retsu Kaioh"           "retsu@shinshinkai.example"       "$PASS"
register "Katsumi Orochi"        "katsumi@shinshinkai.example"     "$PASS"
register "Jack Hanma"            "jack.hanma@example.example"              "$PASS"

register "Ikki (Fénix)"          "ikki@saint.example"              "$PASS"

say "login Roshi (debería dar 200)"
login "roshi@kamehouse.example" "$PASS"
