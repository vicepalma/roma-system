#!/usr/bin/env bash
set -euo pipefail

# Requisitos:
#  - backend corriendo en http://localhost:8080
#  - contenedor postgres llamado "roma-postgres"
#  - jq instalado
#  - usuarios/programas ya seedados (como hicimos)

API=http://localhost:8080
PSQL="docker exec -i roma-postgres psql -U roma -d roma -t -A -c"

today() { date +%F; }

login() {
  local email="$1" pass="$2"
  curl -s -X POST "$API/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"email\":\"$email\",\"password\":\"$pass\"}" \
    | jq -r '.tokens.access'
}

post_json_auth() {
  local token="$1" path="$2" body="$3"
  curl -s -X POST "$API$path" \
    -H "Authorization: Bearer $token" \
    -H 'Content-Type: application/json' \
    -d "$body"
}

echo "== 1) Logins =="
ACCESS_BAKI=$(login "baki.hanma@example.example" "secret123")
ACCESS_ROSHI=$(login "roshi@kamehouse.example" "secret123")
ACCESS_IKKI=$(login "ikki@saint.example" "secret123")

test "${#ACCESS_BAKI}" -gt 20 || { echo "Login Baki falló"; exit 1; }
test "${#ACCESS_ROSHI}" -gt 20 || { echo "Login Roshi falló"; exit 1; }
test "${#ACCESS_IKKI}" -gt 20 || { echo "Login Ikki falló"; exit 1; }
echo "   OK logins"

echo "== 2) IDs de usuarios == "
BAKI_ID=$($PSQL "SELECT id FROM users WHERE email='baki.hanma@example.example';"); echo "BAKI_ID=$BAKI_ID"
RETSU_ID=$($PSQL "SELECT id FROM users WHERE email='retsu@shinshinkai.example';"); echo "RETSU_ID=$RETSU_ID"
KATSUMI_ID=$($PSQL "SELECT id FROM users WHERE email='katsumi@shinshinkai.example';"); echo "KATSUMI_ID=$KATSUMI_ID"
JACK_ID=$($PSQL "SELECT id FROM users WHERE email='jack.hanma@example.example';"); echo "JACK_ID=$JACK_ID"

ROSHI_ID=$($PSQL "SELECT id FROM users WHERE email='roshi@kamehouse.example';"); echo "ROSHI_ID=$ROSHI_ID"
KRILLIN_ID=$($PSQL "SELECT id FROM users WHERE email='krillin@kamehouse.example';"); echo "KRILLIN_ID=$KRILLIN_ID"
YAMCHA_ID=$($PSQL "SELECT id FROM users WHERE email='yamcha@capsule.example';"); echo "YAMCHA_ID=$YAMCHA_ID"
GOKU_ID=$($PSQL "SELECT id FROM users WHERE email='goku@capsule.example';"); echo "GOKU_ID=$GOKU_ID"

IKKI_ID=$($PSQL "SELECT id FROM users WHERE email='ikki@saint.example';"); echo "IKKI_ID=$IKKI_ID"

echo "== 3) Programas de cada coach =="
PROG_BAKI=$($PSQL "SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id WHERE u.email='baki.hanma@example.example' ORDER BY p.created_at DESC LIMIT 1;"); echo "PROG_BAKI=$PROG_BAKI"
PROG_ROSHI=$($PSQL "SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id WHERE u.email='roshi@kamehouse.example' ORDER BY p.created_at DESC LIMIT 1;"); echo "PROG_ROSHI=$PROG_ROSHI"
PROG_IKKI=$($PSQL "SELECT p.id FROM programs p JOIN users u ON u.id=p.owner_id WHERE u.email='ikki@saint.example' ORDER BY p.created_at DESC LIMIT 1;"); echo "PROG_IKKI=$PROG_IKKI"

# Helper para crear link (coach -> disciple) con auto_accept
create_link() {
  local token="$1" disciple_id="$2"
  post_json_auth "$token" "/api/coach/links" "{\"disciple_id\":\"$disciple_id\",\"auto_accept\":true}" \
    | jq -r '.id // .error // "ok"'
}

# Helper para asignar programa (coach -> disciple)
assign_program() {
  local token="$1" disciple_id="$2" program_id="$3" start="$(today)"
  post_json_auth "$token" "/api/coach/assignments" \
    "{\"disciple_id\":\"$disciple_id\",\"program_id\":\"$program_id\",\"start_date\":\"$start\"}" \
    | jq -r '.id // .error // "ok"'
}

echo "== 4) Vínculos (coach_links) =="

echo "Baki -> Retsu:   $(create_link "$ACCESS_BAKI" "$RETSU_ID")"
echo "Baki -> Katsumi: $(create_link "$ACCESS_BAKI" "$KATSUMI_ID")"
echo "Baki -> Jack:    $(create_link "$ACCESS_BAKI" "$JACK_ID")"

echo "Roshi -> Krillin: $(create_link "$ACCESS_ROSHI" "$KRILLIN_ID")"
echo "Roshi -> Yamcha:  $(create_link "$ACCESS_ROSHI" "$YAMCHA_ID")"
echo "Roshi -> Goku:    $(create_link "$ACCESS_ROSHI" "$GOKU_ID")"

echo "Ikki -> Ikki (self-coach): $(create_link "$ACCESS_IKKI" "$IKKI_ID")"

echo "== 5) Assignments =="

echo "Baki -> Retsu (prog Baki):     $(assign_program "$ACCESS_BAKI" "$RETSU_ID"   "$PROG_BAKI")"
echo "Baki -> Katsumi (prog Baki):   $(assign_program "$ACCESS_BAKI" "$KATSUMI_ID" "$PROG_BAKI")"
echo "Baki -> Jack (prog Baki):      $(assign_program "$ACCESS_BAKI" "$JACK_ID"    "$PROG_BAKI")"

echo "Roshi -> Krillin (prog Roshi): $(assign_program "$ACCESS_ROSHI" "$KRILLIN_ID" "$PROG_ROSHI")"
echo "Roshi -> Yamcha (prog Roshi):  $(assign_program "$ACCESS_ROSHI" "$YAMCHA_ID"  "$PROG_ROSHI")"
echo "Roshi -> Goku (prog Roshi):    $(assign_program "$ACCESS_ROSHI" "$GOKU_ID"    "$PROG_ROSHI")"

echo "Ikki -> Ikki (prog Ikki):      $(assign_program "$ACCESS_IKKI" "$IKKI_ID"     "$PROG_IKKI")"

echo "== 6) Verificaciones rápidas =="
echo "- Discipulos de Baki:"
curl -s "$API/api/coach/disciples" -H "Authorization: Bearer $ACCESS_BAKI" | jq

echo "- Overview de Retsu (7 días, volume, America/Santiago):"
curl -s "$API/api/coach/disciples/$RETSU_ID/overview?days=7&metric=volume&tz=America/Santiago" \
  -H "Authorization: Bearer $ACCESS_BAKI" | jq
