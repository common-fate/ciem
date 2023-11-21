permit(
    principal == User::"chris",
    action == Action::"view",
    resource in Folder::"prod"
);