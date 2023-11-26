permit (
    principal,
    action == JIT::Action::"Review",
    resource
);
// e.g. policy can be added here to prevent self-reviewing.
// forbid(
//     principal,
//     action == JIT::Action::"Review",
//     resource
// ) when {
//     resource.principal == principal
// };
