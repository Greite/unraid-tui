package api

const querySystemInfo = `query {
	info {
		cpu {
			manufacturer
			brand
			cores
			threads
			packages {
				totalPower
				temp
			}
		}
		os {
			platform
			distro
			release
			uptime
			hostname
			kernel
		}
	}
}`

const querySystemMetrics = `query {
	metrics {
		cpu {
			percentTotal
			cpus {
				percentTotal
			}
		}
		memory {
			used
			total
			available
			percentTotal
		}
	}
}`

const queryNotificationsOverview = `query {
	notifications {
		overview {
			unread {
				info
				warning
				alert
				total
			}
		}
	}
}`

const queryNotificationsList = `query {
	notifications {
		list(filter: { type: UNREAD, offset: 0, limit: 100 }) {
			id
			title
			subject
			description
			importance
			timestamp
		}
	}
}`

const queryNetwork = `query {
	network {
		accessUrls {
			name
			type
			ipv4
			ipv6
		}
	}
}`

const queryShares = `query {
	shares {
		name
		free
		used
		size
		cache
		comment
	}
}`

const queryArrayState = `query {
	array {
		state
		capacity {
			kilobytes {
				free
				used
				total
			}
		}
		parityCheckStatus {
			status
			progress
			running
		}
	}
}`

const queryDisks = `query {
	array {
		disks {
			name device size fsSize fsFree fsUsed status type temp
		}
		caches {
			name device size fsSize fsFree fsUsed status type temp
		}
		parities {
			name device size fsSize fsFree fsUsed status type temp
		}
	}
}`

// Container mutation templates — %s is replaced with the container ID.
const mutationStartContainer = `mutation { docker { start(id: "%s") { id state } } }`
const mutationStopContainer = `mutation { docker { stop(id: "%s") { id state } } }`
const mutationPauseContainer = `mutation { docker { pause(id: "%s") { id state } } }`
const mutationUnpauseContainer = `mutation { docker { unpause(id: "%s") { id state } } }`
const mutationUpdateContainer = `mutation { docker { updateContainer(id: "%s") { id names image } } }`
const mutationUpdateAllContainers = `mutation { docker { updateAllContainers { id names } } }`

const queryVMs = `query {
	vms {
		domains {
			id
			name
			state
		}
	}
}`

const queryContainers = `query {
	docker {
		containers {
			id
			names
			image
			state
			status
			autoStart
			ports {
				privatePort
				publicPort
				type
			}
			webUiUrl
			isUpdateAvailable
		}
	}
}`
