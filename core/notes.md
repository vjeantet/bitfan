# Modify pipeline on the fly
Pipeline's agents can be updated and reloaded (pause/resume)
* reload = pause and resume with a updated configuration
* the inner agent's processor stop and start when agent pause and resume

Pipeline's connections can not be updated ! when needed, the entire pipeline sould be restarted
* a connection is defined by a source, a destination and a buffer size
