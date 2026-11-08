# Logical lineages

Use these rules to distinguish independent recurrence from retries and related records.

| Evidence | Treatment |
| --- | --- |
| The same `source_ref`, including verified versions after resume | One logical source. |
| Explicit `parent` or `forked_from` links, including transitive links | One lineage. |
| Subagent evidence | Include in the parent's lineage only when an explicit relationship or normalized parent association establishes it. |
| Different references within one source instance | Count as independent only when relationship coverage is complete and supplied context leaves no ambiguity. |
| Different source instances | Count as independent only when supplied evidence establishes separate interactions, even if each response is complete. |
| Same-provider sources in different instances with the same provider-source fingerprint | Leave unresolved unless evidence distinguishes copied observations from independent interactions. |

Build only the smallest relationship graph the evidence supports.\
Do not merge records merely because their text or failures are similar.\
When a required relationship is missing or ambiguous, preserve that uncertainty and do not count the affected sources as confirmed independent recurrence.

Report retries within each lineage and the number of independent lineages separately.
