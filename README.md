# f-hsub
Find hsub messages in an mbox file.
This software was created for the program suck,
which fetches all (new) articles from Usenet Groups.

When valid hsubs are found, in an mbox file, each article will be seperately saved.

[2024-23-06]
Added a feature to look not only for a hashed 'Subject:' but also
for a 'H-Hsub:' hashed header, in case people like to keep the
orignal Subject line, without hsub content.
