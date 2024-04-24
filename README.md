## Task monster rpg

### NEW

Concept:

1. Each user have an avatar
2. The user sets a goal - to learn a new language, complete a project, etc.
3. Goal was autofilled with chatGPT and similar content
4. The more user achieve in that goal - the more he is stronger against some monsters
5. Monster fights with users and each other in real time
6. Each victory - unique info/links/quests for user's goal

### OLD

Concept:

1. The user sets a goal - to learn a new language, complete a project, etc.

2. Key features are specified in the goal - deadlines, sources of information, useful links, etc.

3. Some "concepts" are parsed from useful links (using ChatGPT)

4. "Concepts" are presented to the user, who must either write something based on them, memorize them, or perform an action

5. All of this is gamified, tasks are like raids, actions give experience points, who will be the monster is still a question

Differences from HabitRPG - the presence of a parser for collecting information/formatting it into tasks and motivating the user to complete them.
The point is that there is always project information that the user has missed. In some goals, it can be difficult to figure out where to move next.

Contact with the user through a bot, currently via Telegram.

API operation scheme: - set the topic by token (userID or other) and receive a list of related topics

###todo
- Storing data in a raised database
- Integration with a Telegram bot
- Makefile for deploying the service
- Move everything to Docker
- Create a history with visualization for the player's story
