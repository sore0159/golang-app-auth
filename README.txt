Authorization and Userdata app for project server.

Designed to handle user profiles, authorization, and general request context.  Provides wrappers for other apps to use in building their own context solutions.

Allows for tiered access to projects (handled by the projects themselves) and provides directories for users to save appdata (Using gob at present).

Sets up logging for users that copies user-specific log events to a separate file for each user (one file for all guests).

Allows guests to have thier own app data that they can access until they close their browser or midnight.  50 guest data dirs per day current limit: no dir is created until a guest actually needs one.

Passwords are saved as plain text in files.  This is bad and I feel bad.
