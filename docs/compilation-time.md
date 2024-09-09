# Compilation Time
When using our automatic instrumentation tool,
two additional phases are added before the above steps: **Preprocessing** and **Instrument**.
The total compilation time is the sum of the time taken by these two phases and
the time taken by the original compilation process. In general, `~92.8%` of the 
total compilation time is increased due to these two phases.