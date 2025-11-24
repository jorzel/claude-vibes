#!/bin/bash

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
DURATION="${DURATION:-60}"  # Duration in seconds
REQUESTS_PER_SECOND="${RPS:-5}"  # Requests per second
VERBOSE="${VERBOSE:-false}"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Storage for created resources
declare -a EVENT_IDS
declare -a BOOKING_IDS
USER_ID="$(uuidgen)"

log() {
    if [ "$VERBOSE" = "true" ]; then
        echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
    fi
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Check if server is running
check_server() {
    if ! curl -s -f "${BASE_URL}/health" > /dev/null 2>&1; then
        error "Server is not running at ${BASE_URL}"
        error "Start your server first!"
        exit 1
    fi
    echo -e "${GREEN}Server is running${NC}"
}

# Create an event
create_event() {
    local name="$1"
    local date="$2"
    local location="$3"
    local tickets="$4"

    response=$(curl -s -X POST "${BASE_URL}/events" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"${name}\",\"date\":\"${date}\",\"location\":\"${location}\",\"tickets\":${tickets}}")

    event_id=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$event_id" ]; then
        EVENT_IDS+=("$event_id")
        log "Created event: $name (ID: $event_id)"
    fi
}

# Get list of events
list_events() {
    curl -s "${BASE_URL}/events" > /dev/null
    log "Listed all events"
}

# Get specific event
get_event() {
    if [ ${#EVENT_IDS[@]} -gt 0 ]; then
        local random_idx=$((RANDOM % ${#EVENT_IDS[@]}))
        local event_id="${EVENT_IDS[$random_idx]}"
        curl -s "${BASE_URL}/events/${event_id}" > /dev/null
        log "Retrieved event: $event_id"
    fi
}

# Create a booking
create_booking() {
    if [ ${#EVENT_IDS[@]} -eq 0 ]; then
        return
    fi

    local random_idx=$((RANDOM % ${#EVENT_IDS[@]}))
    local event_id="${EVENT_IDS[$random_idx]}"
    local tickets=$((1 + RANDOM % 5))

    response=$(curl -s -X POST "${BASE_URL}/bookings" \
        -H "Content-Type: application/json" \
        -d "{\"event_id\":\"${event_id}\",\"user_id\":\"${USER_ID}\",\"tickets_booked\":${tickets}}")

    booking_id=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$booking_id" ]; then
        BOOKING_IDS+=("$booking_id")
        log "Created booking: $booking_id for event $event_id"
    fi
}

# Get specific booking
get_booking() {
    if [ ${#BOOKING_IDS[@]} -gt 0 ]; then
        local random_idx=$((RANDOM % ${#BOOKING_IDS[@]}))
        local booking_id="${BOOKING_IDS[$random_idx]}"
        curl -s "${BASE_URL}/bookings/${booking_id}" > /dev/null
        log "Retrieved booking: $booking_id"
    fi
}

# Health check
health_check() {
    curl -s "${BASE_URL}/health" > /dev/null
    log "Health check"
}

# Random request generator
random_request() {
    local action=$((RANDOM % 100))

    if [ "$action" -lt 5 ]; then
        # 5% - Create event
        local event_num=$((RANDOM % 1000))
        local future_date=$(date -u -v+30d +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "+30 days" +"%Y-%m-%dT%H:%M:%SZ")
        create_event "Event-${event_num}" "$future_date" "Location-${event_num}" "$((50 + RANDOM % 200))"
    elif [ "$action" -lt 15 ]; then
        # 10% - List events
        list_events
    elif [ "$action" -lt 35 ]; then
        # 20% - Get event
        get_event
    elif [ "$action" -lt 60 ]; then
        # 25% - Create booking
        create_booking
    elif [ "$action" -lt 75 ]; then
        # 15% - Get booking
        get_booking
    else
        # 25% - Health check
        health_check
    fi
}

# Initialize with some events
initialize_data() {
    echo -e "${YELLOW}Initializing test data...${NC}"
    local future_date=$(date -u -v+30d +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "+30 days" +"%Y-%m-%dT%H:%M:%SZ")

    create_event "Concert-2025" "$future_date" "Madison Square Garden" 1000
    create_event "Tech Conference" "$future_date" "Silicon Valley" 500
    create_event "Sports Game" "$future_date" "Stadium" 2000
    create_event "Theater Show" "$future_date" "Broadway" 300
    create_event "Festival" "$future_date" "Central Park" 5000

    echo -e "${GREEN}Created ${#EVENT_IDS[@]} initial events${NC}"
}

# Main traffic simulation loop
simulate_traffic() {
    echo -e "${YELLOW}Starting traffic simulation...${NC}"
    echo "Duration: ${DURATION}s"
    echo "Rate: ${REQUESTS_PER_SECOND} req/s"
    echo "Target: ${BASE_URL}"
    echo ""

    local end_time=$(($(date +%s) + DURATION))
    local request_count=0
    local sleep_time=$(echo "scale=4; 1.0 / $REQUESTS_PER_SECOND" | bc)

    while [ $(date +%s) -lt $end_time ]; do
        random_request
        request_count=$((request_count + 1))

        # Print progress every 10 requests
        if [ $((request_count % 10)) -eq 0 ]; then
            local elapsed=$(($(date +%s) - (end_time - DURATION)))
            echo -ne "\rRequests sent: ${request_count} | Elapsed: ${elapsed}s | Events: ${#EVENT_IDS[@]} | Bookings: ${#BOOKING_IDS[@]}   "
        fi

        sleep "$sleep_time" 2>/dev/null || sleep 0.2
    done

    echo ""
    echo -e "${GREEN}Simulation complete!${NC}"
    echo "Total requests: ${request_count}"
    echo "Events created: ${#EVENT_IDS[@]}"
    echo "Bookings created: ${#BOOKING_IDS[@]}"
    echo ""
    echo "View metrics at: ${BASE_URL}/metrics"
}

# Print usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Simulate HTTP traffic to the booking service for metrics collection.

OPTIONS:
    -h, --help              Show this help message
    -u, --url URL           Base URL (default: http://localhost:8080)
    -d, --duration SECONDS  Duration in seconds (default: 60)
    -r, --rps RATE          Requests per second (default: 5)
    -v, --verbose           Enable verbose logging

EXAMPLES:
    # Run for 60 seconds at 5 req/s
    $0

    # Run for 2 minutes at 10 req/s
    $0 -d 120 -r 10

    # Run with verbose logging
    $0 -v

    # Custom server URL
    $0 -u http://localhost:9000

ENVIRONMENT VARIABLES:
    BASE_URL    Server base URL
    DURATION    Duration in seconds
    RPS         Requests per second
    VERBOSE     Enable verbose output (true/false)

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -u|--url)
            BASE_URL="$2"
            shift 2
            ;;
        -d|--duration)
            DURATION="$2"
            shift 2
            ;;
        -r|--rps)
            REQUESTS_PER_SECOND="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE="true"
            shift
            ;;
        *)
            error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Main execution
echo -e "${YELLOW}=== Booking Service Traffic Simulator ===${NC}\n"

check_server
initialize_data
simulate_traffic

echo -e "\n${GREEN}Done! Check your metrics at ${BASE_URL}/metrics${NC}"
