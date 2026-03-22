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

const querySystemInfoExtra = `query {
	info {
		versions {
			core {
				unraid
				api
				kernel
			}
		}
		devices {
			gpu { id name vendor model }
			pci { id name vendor model }
			usb { id name vendor model }
		}
		memory {
			layout {
				size
				type
				clockSpeed
				manufacturer
				bank
			}
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

const queryParityHistory = `query {
	array {
		parityHistory {
			date
			status
			duration
			speed
			errors
		}
	}
}`

// Docker queries
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
		}
		containerUpdateStatuses {
			name
			updateStatus
		}
	}
}`

const queryContainerStats = `query {
	docker {
		containers {
			id
			names
			state
			cpuPercent
			memUsage
			memPercent
		}
	}
}`

// Docker mutations — %s is replaced with the container ID.
const mutationStartContainer = `mutation { docker { start(id: "%s") { id state } } }`
const mutationStopContainer = `mutation { docker { stop(id: "%s") { id state } } }`
const mutationPauseContainer = `mutation { docker { pause(id: "%s") { id state } } }`
const mutationUnpauseContainer = `mutation { docker { unpause(id: "%s") { id state } } }`
const mutationUpdateContainer = `mutation { docker { updateContainer(id: "%s") { id names image } } }`
const mutationUpdateAllContainers = `mutation { docker { updateAllContainers { id names } } }`
const mutationAutostart = `mutation { docker { updateAutostartConfiguration(entries: [{ id: "%s", autoStart: %t, wait: %d }]) } }`

// VM queries
const queryVMs = `query {
	vms {
		domains {
			id
			name
			state
		}
	}
}`

// VM mutations
const mutationVMStart = `mutation { vm { start(id: "%s") } }`
const mutationVMStop = `mutation { vm { stop(id: "%s") } }`
const mutationVMPause = `mutation { vm { pause(id: "%s") } }`
const mutationVMResume = `mutation { vm { resume(id: "%s") } }`
const mutationVMForceStop = `mutation { vm { forceStop(id: "%s") } }`
const mutationVMReboot = `mutation { vm { reboot(id: "%s") } }`

// Notification mutations
const mutationArchiveNotification = `mutation { archiveNotification(id: "%s") { id } }`
const mutationArchiveAllNotifications = `mutation { archiveAll { unread { total } } }`
