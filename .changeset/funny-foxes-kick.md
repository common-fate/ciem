---
"@common-fate/cli": patch
---

If a single policy is being validated using `cf authz policyset validate`, the CLI will print "validating 1 policy" rather than "validating 1 policies".
