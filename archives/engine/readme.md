Role: You are a Senior Backend Engineer specializing in High-Performance Go Systems and Distributed Architecture.

This engine is heart of Business rule web application.

Task: Design a "Micro-Engine Logic Router" in Go that processes incoming JSON data by routing it to specific, independent Go API services (Rules) based on URL parameters. 
Constraints:
 * No External Libraries: Use only the Go Standard Library (net/http, encoding/json, net/http/httputil, etc.).
 * Architecture: * The Orchestrator: A central gateway that accepts JSON payloads.
   * The Registry: A dynamic mapping of rule_names to service_urls.
   * The Worker Template: A lightweight Go boilerplate for the individual "Rule Services."
 * Hot-Reloading: The Orchestrator must be able to update its routing registry at runtime (using a sync.RWMutex or atomic.Value) without a restart.
 * Efficiency: Use httputil.ReverseProxy or optimized io.Copy for zero-copy-style data forwarding to minimize memory overhead.
 * Complex Logic: Provide an example within a Rule Service that handles nested JSON looping and conditional transformations.
Expected Output Structure:
 * orchestrator/main.go: Featuring the dynamic router and a "Refresh Registry" endpoint.
 * rules/template.go: A reusable base for creating new logic services.
 * Internal Communication Design: A strategy for how these services should handle timeouts and health checks.
 * Deployment Guide: A brief explanation of how to "Hot Swap" a rule service without dropping traffic.
Why use this prompt?
 * Focuses on ReverseProxy: In Go, the net/http/httputil package has a built-in ReverseProxy. Using this is much more efficient than manually reading and re-writing the JSON body in the Orchestrator, as it streams the data directly.
 * Thread Safety: By mentioning sync.RWMutex, you ensure the AI writes code that won't crash when you update a rule's URL while traffic is flowing.
 * Standard Library: It prevents the AI from suggesting "heavy" frameworks like Gin or Echo, keeping your "in-house" requirement intact.