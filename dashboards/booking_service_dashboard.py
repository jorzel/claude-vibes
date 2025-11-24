#!/usr/bin/env python3
"""
Booking Service Grafana Dashboard Generator

This script generates a Grafana dashboard for monitoring the Booking Service.
It uses grafanalib to define the dashboard as code.

Usage:
    python dashboards/booking_service_dashboard.py > dashboards/booking_service.json

Requirements:
    pip install grafanalib
"""

from grafanalib.core import (
    Dashboard,
    Graph,
    Row,
    Target,
    Templating,
    Template,
    SingleStat,
    YAxes,
    YAxis,
    SECONDS_FORMAT,
    SHORT_FORMAT,
    OPS_FORMAT,
)


def create_dashboard():
    """Create the Booking Service monitoring dashboard."""

    return Dashboard(
        title="Booking Service Monitoring",
        description="Comprehensive monitoring dashboard for the Booking Service",
        tags=["booking-service", "events", "bookings"],
        timezone="browser",
        refresh="30s",
        templating=Templating(
            list=[
                Template(
                    name="datasource",
                    label="Data Source",
                    type="datasource",
                    query="prometheus",
                ),
            ]
        ),
        rows=[
            # Row 1: Overview Stats
            Row(
                title="Overview",
                panels=[
                    SingleStat(
                        title="Total Events Created",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='sum(booking_service_events_created_total)',
                                legendFormat="Events",
                            ),
                        ],
                        valueName="current",
                        format=SHORT_FORMAT,
                        sparklineShow=True,
                    ),
                    SingleStat(
                        title="Total Bookings Created",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='sum(booking_service_bookings_created_total)',
                                legendFormat="Bookings",
                            ),
                        ],
                        valueName="current",
                        format=SHORT_FORMAT,
                        sparklineShow=True,
                    ),
                    SingleStat(
                        title="Total Tickets Booked",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='sum(booking_service_tickets_booked_total)',
                                legendFormat="Tickets",
                            ),
                        ],
                        valueName="current",
                        format=SHORT_FORMAT,
                        sparklineShow=True,
                    ),
                    SingleStat(
                        title="Request Rate (req/s)",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='rate(booking_service_http_request_duration_seconds_count[5m])',
                                legendFormat="Rate",
                            ),
                        ],
                        valueName="current",
                        format=OPS_FORMAT,
                        sparklineShow=True,
                    ),
                ],
            ),

            # Row 2: Business Metrics
            Row(
                title="Business Metrics",
                panels=[
                    Graph(
                        title="Events Created Over Time",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='rate(booking_service_events_created_total{status="success"}[5m])',
                                legendFormat="Success",
                            ),
                            Target(
                                expr='rate(booking_service_events_created_total{status="error"}[5m])',
                                legendFormat="Error",
                            ),
                        ],
                        yAxes=YAxes(
                            left=YAxis(format=OPS_FORMAT),
                        ),
                    ),
                    Graph(
                        title="Bookings Created Over Time",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='rate(booking_service_bookings_created_total{status="success"}[5m])',
                                legendFormat="Success",
                            ),
                            Target(
                                expr='rate(booking_service_bookings_created_total{status="error"}[5m])',
                                legendFormat="Error",
                            ),
                        ],
                        yAxes=YAxes(
                            left=YAxis(format=OPS_FORMAT),
                        ),
                    ),
                ],
            ),

            # Row 3: HTTP Performance
            Row(
                title="HTTP Performance",
                panels=[
                    Graph(
                        title="Request Rate by Endpoint",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='rate(booking_service_http_request_duration_seconds_count[5m])',
                                legendFormat="{{method}} {{path}}",
                            ),
                        ],
                        yAxes=YAxes(
                            left=YAxis(format=OPS_FORMAT),
                        ),
                    ),
                    Graph(
                        title="Response Time (p50, p95, p99)",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='histogram_quantile(0.50, rate(booking_service_http_request_duration_seconds_bucket[5m]))',
                                legendFormat="p50",
                            ),
                            Target(
                                expr='histogram_quantile(0.95, rate(booking_service_http_request_duration_seconds_bucket[5m]))',
                                legendFormat="p95",
                            ),
                            Target(
                                expr='histogram_quantile(0.99, rate(booking_service_http_request_duration_seconds_bucket[5m]))',
                                legendFormat="p99",
                            ),
                        ],
                        yAxes=YAxes(
                            left=YAxis(format=SECONDS_FORMAT),
                        ),
                    ),
                ],
            ),

            # Row 4: Error Rates
            Row(
                title="Error Rates",
                panels=[
                    Graph(
                        title="HTTP Error Rate by Status Code",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='rate(booking_service_http_request_duration_seconds_count{status=~"4.."}[5m])',
                                legendFormat="4xx - {{status}}",
                            ),
                            Target(
                                expr='rate(booking_service_http_request_duration_seconds_count{status=~"5.."}[5m])',
                                legendFormat="5xx - {{status}}",
                            ),
                        ],
                        yAxes=YAxes(
                            left=YAxis(format=OPS_FORMAT),
                        ),
                    ),
                    Graph(
                        title="Error Rate Percentage",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='rate(booking_service_http_request_duration_seconds_count{status=~"[45].."}[5m]) / rate(booking_service_http_request_duration_seconds_count[5m]) * 100',
                                legendFormat="Error %",
                            ),
                        ],
                        yAxes=YAxes(
                            left=YAxis(format="percent", min=0, max=100),
                        ),
                    ),
                ],
            ),

            # Row 5: System Metrics
            Row(
                title="System Metrics",
                panels=[
                    Graph(
                        title="Go Goroutines",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='go_goroutines',
                                legendFormat="Goroutines",
                            ),
                        ],
                        yAxes=YAxes(
                            left=YAxis(format=SHORT_FORMAT),
                        ),
                    ),
                    Graph(
                        title="Memory Usage",
                        dataSource="$datasource",
                        targets=[
                            Target(
                                expr='go_memstats_alloc_bytes',
                                legendFormat="Allocated",
                            ),
                            Target(
                                expr='go_memstats_heap_alloc_bytes',
                                legendFormat="Heap",
                            ),
                        ],
                        yAxes=YAxes(
                            left=YAxis(format="bytes"),
                        ),
                    ),
                ],
            ),
        ],
    ).auto_panel_ids()


if __name__ == "__main__":
    import json
    from grafanalib.core import _gen_json_from_panels

    dashboard = create_dashboard()

    # Generate JSON
    dashboard_json = {
        "dashboard": json.loads(_gen_json_from_panels(dashboard)),
        "overwrite": True,
        "inputs": [],
    }

    print(json.dumps(dashboard_json, indent=2))
