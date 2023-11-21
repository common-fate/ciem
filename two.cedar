permit(
    principal == User::"josh",
    action == Action::"view",
    resource in Folder::"test"
);