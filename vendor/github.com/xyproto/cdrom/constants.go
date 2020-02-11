package cdrom

import (
	"math"
)

const (
	EDRIVE_CANT_DO_THIS = 95
	EOPNOTSUPP          = 95

	/*******************************************************
	 * The CD-ROM IOCTL commands  -- these should be supported by
	 * all the various cdrom drivers.  For the CD-ROM ioctls, we
	 * will commandeer byte 0x53, or 'S'.
	 *******************************************************/
	CDROMPAUSE      = 0x5301 /* Pause Audio Operation */
	CDROMRESUME     = 0x5302 /* Resume paused Audio Operation */
	CDROMPLAYMSF    = 0x5303 /* Play Audio MSF (struct cdrom_msf) */
	CDROMPLAYTRKIND = 0x5304 /* Play Audio Track/index
	   (struct cdrom_ti) */
	CDROMREADTOCHDR = 0x5305 /* Read TOC header
	   (struct cdrom_tochdr) */
	CDROMREADTOCENTRY = 0x5306 /* Read TOC entry
	   (struct cdrom_tocentry) */
	CDROMSTOP    = 0x5307 /* Stop the cdrom drive */
	CDROMSTART   = 0x5308 /* Start the cdrom drive */
	CDROMEJECT   = 0x5309 /* Ejects the cdrom media */
	CDROMVOLCTRL = 0x530a /* Control output volume
	   (struct cdrom_volctrl) */
	CDROMSUBCHNL = 0x530b /* Read subchannel data
	   (struct cdrom_subchnl) */
	CDROMREADMODE2 = 0x530c /* Read CDROM mode 2 data (2336 Bytes)
	   (struct cdrom_read) */
	CDROMREADMODE1 = 0x530d /* Read CDROM mode 1 data (2048 Bytes)
	   (struct cdrom_read) */
	CDROMREADAUDIO    = 0x530e /* (struct cdrom_read_audio) */
	CDROMEJECT_SW     = 0x530f /* enable(1)/disable(0) auto-ejecting */
	CDROMMULTISESSION = 0x5310 /* Obtain the start-of-last-session
	   address of multi session disks
	   (struct cdrom_multisession) */
	CDROM_GET_MCN = 0x5311 /* Obtain the "Universal Product Code"
	   if available (struct cdrom_mcn) */
	CDROM_GET_UPC = CDROM_GET_MCN /* This one is deprecated,
	   but here anyway for compatibility */
	CDROMRESET   = 0x5312 /* hard-reset the drive */
	CDROMVOLREAD = 0x5313 /* Get the drive's volume setting
	   (struct cdrom_volctrl) */
	CDROMREADRAW = 0x5314 /* read data in raw mode (2352 Bytes)
	   (struct cdrom_read) */
	/*
	 * These ioctls are used only used in aztcd.c and optcd.c
	 */
	CDROMREADCOOKED = 0x5315 /* read data in cooked mode */
	CDROMSEEK       = 0x5316 /* seek msf address */

	/*
	   * This ioctl is only used by the scsi-cd driver.
	     It is for playing audio in logical block addressing mode.
	*/
	CDROMPLAYBLK = 0x5317 /* (struct cdrom_blk) */

	/*
	 * These ioctls are only used in optcd.c
	 */
	CDROMREADALL = 0x5318 /* read all 2646 bytes */

	/*
	 * These ioctls are (now) only in ide-cd.c for controlling
	 * drive spindown time.  They should be implemented in the
	 * Uniform driver, via generic packet commands, GPCMD_MODE_SELECT_10,
	 * GPCMD_MODE_SENSE_10 and the GPMODE_POWER_PAGE...
	 *  -Erik
	 */
	CDROMGETSPINDOWN = 0x531d
	CDROMSETSPINDOWN = 0x531e

	/*
	 * These ioctls are implemented through the uniform CD-ROM driver
	 * They _will_ be adopted by all CD-ROM drivers, when all the CD-ROM
	 * drivers are eventually ported to the uniform CD-ROM driver interface.
	 */
	CDROMCLOSETRAY       = 0x5319 /* pendant of CDROMEJECT */
	CDROM_SET_OPTIONS    = 0x5320 /* Set behavior options */
	CDROM_CLEAR_OPTIONS  = 0x5321 /* Clear behavior options */
	CDROM_SELECT_SPEED   = 0x5322 /* Set the CD-ROM speed */
	CDROM_SELECT_DISC    = 0x5323 /* Select disc (for juke-boxes) */
	CDROM_MEDIA_CHANGED  = 0x5325 /* Check is media changed  */
	CDROM_DRIVE_STATUS   = 0x5326 /* Get tray position, etc. */
	CDROM_DISC_STATUS    = 0x5327 /* Get disc type, etc. */
	CDROM_CHANGER_NSLOTS = 0x5328 /* Get number of slots */
	CDROM_LOCKDOOR       = 0x5329 /* lock or unlock door */
	CDROM_DEBUG          = 0x5330 /* Turn debug messages on/off */
	CDROM_GET_CAPABILITY = 0x5331 /* get capabilities */

	/* Note that scsi/scsi_ioctl.h also uses 0x5382 - 0x5386.
	 * Future CDROM ioctls should be kept below 0x537F
	 */

	/* This ioctl is only used by sbpcd at the moment */
	CDROMAUDIOBUFSIZ = 0x5382 /* set the audio buffer size */
	/* conflict with SCSI_IOCTL_GET_IDLUN */

	/* DVD-ROM Specific ioctls */
	DVD_READ_STRUCT  = 0x5390 /* Read structure */
	DVD_WRITE_STRUCT = 0x5391 /* Write structure */
	DVD_AUTH         = 0x5392 /* Authentication */

	CDROM_SEND_PACKET   = 0x5393 /* send a packet to the drive */
	CDROM_NEXT_WRITABLE = 0x5394 /* get next writable block */
	CDROM_LAST_WRITTEN  = 0x5395 /* get last block written on disc */

	CDROM_PACKET_SIZE = 12

	CGC_DATA_UNKNOWN = 0
	CGC_DATA_WRITE   = 1
	CGC_DATA_READ    = 2
	CGC_DATA_NONE    = 3

	/*
	    * A CD-ROM physical sector size is 2048, 2052, 2056, 2324, 2332, 2336,
	    * 2340, or 2352 bytes long.

	   *         Sector types of the standard CD-ROM data formats:
	    *
	    * format   sector type               user data size (bytes)
	    * -----------------------------------------------------------------------------
	    *   1     (Red Book)    CD-DA          2352    (CD_FRAMESIZE_RAW)
	    *   2     (Yellow Book) Mode1 Form1    2048    (CD_FRAMESIZE)
	    *   3     (Yellow Book) Mode1 Form2    2336    (CD_FRAMESIZE_RAW0)
	    *   4     (Green Book)  Mode2 Form1    2048    (CD_FRAMESIZE)
	    *   5     (Green Book)  Mode2 Form2    2328    (2324+4 spare bytes)
	    *
	    *
	    *       The layout of the standard CD-ROM data formats:
	    * -----------------------------------------------------------------------------
	    * - audio (red):                  | audio_sample_bytes |
	    *                                 |        2352        |
	    *
	    * - data (yellow, mode1):         | sync - head - data - EDC - zero - ECC |
	    *                                 |  12  -   4  - 2048 -  4  -   8  - 276 |
	    *
	    * - data (yellow, mode2):         | sync - head - data |
	    *                                 |  12  -   4  - 2336 |
	    *
	    * - XA data (green, mode2 form1): | sync - head - sub - data - EDC - ECC |
	    *                                 |  12  -   4  -  8  - 2048 -  4  - 276 |
	    *
	    * - XA data (green, mode2 form2): | sync - head - sub - data - Spare |
	    *                                 |  12  -   4  -  8  - 2324 -  4    |
	    *
	*/

	/* Some generally useful CD-ROM information -- mostly based on the above */
	CD_MINS            = 74   /* max. minutes per CD, not really a limit */
	CD_SECS            = 60   /* seconds per minute */
	CD_FRAMES          = 75   /* frames per second */
	CD_SYNC_SIZE       = 12   /* 12 sync bytes per raw data frame */
	CD_MSF_OFFSET      = 150  /* MSF numbering offset of first frame */
	CD_CHUNK_SIZE      = 24   /* lowest-level "data bytes piece" */
	CD_NUM_OF_CHUNKS   = 98   /* chunks per frame */
	CD_FRAMESIZE_SUB   = 96   /* subchannel data "frame" size */
	CD_HEAD_SIZE       = 4    /* header (address) bytes per raw data frame */
	CD_SUBHEAD_SIZE    = 8    /* subheader bytes per raw XA data frame */
	CD_EDC_SIZE        = 4    /* bytes EDC per most raw data frame types */
	CD_ZERO_SIZE       = 8    /* bytes zero per yellow book mode 1 frame */
	CD_ECC_SIZE        = 276  /* bytes ECC per most raw data frame types */
	CD_FRAMESIZE       = 2048 /* bytes per frame, "cooked" mode */
	CD_FRAMESIZE_RAW   = 2352 /* bytes per frame, "raw" mode */
	CD_FRAMESIZE_RAWER = 2646 /* The maximum possible returned bytes */
	/* most drives don't deliver everything: */
	CD_FRAMESIZE_RAW1 = (CD_FRAMESIZE_RAW - CD_SYNC_SIZE)                /*2340*/
	CD_FRAMESIZE_RAW0 = (CD_FRAMESIZE_RAW - CD_SYNC_SIZE - CD_HEAD_SIZE) /*2336*/

	CD_XA_HEAD      = (CD_HEAD_SIZE + CD_SUBHEAD_SIZE) /* "before data" part of raw XA frame */
	CD_XA_TAIL      = (CD_EDC_SIZE + CD_ECC_SIZE)      /* "after data" part of raw XA frame */
	CD_XA_SYNC_HEAD = (CD_SYNC_SIZE + CD_XA_HEAD)      /* sync bytes + header of XA frame */

	/* CD-ROM address types (cdrom_tocentry.cdte_format) */
	CDROM_LBA = 0x01 /* "logical block": first frame is #0 */
	CDROM_MSF = 0x02 /* "minute-second-frame": binary, not bcd here! */

	/* bit to tell whether track is data or audio (cdrom_tocentry.cdte_ctrl) */
	CDROM_DATA_TRACK = 0x04

	/* The leadout track is always 0xAA, regardless of # of tracks on disc */
	CDROM_LEADOUT = 0xAA

	/* audio states (from SCSI-2, but seen with other drives, too) */
	CDROM_AUDIO_INVALID   = 0x00 /* audio status not supported */
	CDROM_AUDIO_PLAY      = 0x11 /* audio play operation in progress */
	CDROM_AUDIO_PAUSED    = 0x12 /* audio play operation paused */
	CDROM_AUDIO_COMPLETED = 0x13 /* audio play successfully completed */
	CDROM_AUDIO_ERROR     = 0x14 /* audio play stopped due to error */
	CDROM_AUDIO_NO_STATUS = 0x15 /* no current audio status to return */

	/* capability flags used with the uniform CD-ROM driver */
	CDC_CLOSE_TRAY     = 0x1      /* caddy systems _can't_ close */
	CDC_OPEN_TRAY      = 0x2      /* but _can_ eject.  */
	CDC_LOCK           = 0x4      /* disable manual eject */
	CDC_SELECT_SPEED   = 0x8      /* programmable speed */
	CDC_SELECT_DISC    = 0x10     /* select disc from juke-box */
	CDC_MULTI_SESSION  = 0x20     /* read sessions>1 */
	CDC_MCN            = 0x40     /* Medium Catalog Number */
	CDC_MEDIA_CHANGED  = 0x80     /* media changed */
	CDC_PLAY_AUDIO     = 0x100    /* audio functions */
	CDC_RESET          = 0x200    /* hard reset device */
	CDC_DRIVE_STATUS   = 0x800    /* driver implements drive status */
	CDC_GENERIC_PACKET = 0x1000   /* driver implements generic packets */
	CDC_CD_R           = 0x2000   /* drive is a CD-R */
	CDC_CD_RW          = 0x4000   /* drive is a CD-RW */
	CDC_DVD            = 0x8000   /* drive is a DVD */
	CDC_DVD_R          = 0x10000  /* drive can write DVD-R */
	CDC_DVD_RAM        = 0x20000  /* drive can write DVD-RAM */
	CDC_MO_DRIVE       = 0x40000  /* drive is an MO device */
	CDC_MRW            = 0x80000  /* drive can read MRW */
	CDC_MRW_W          = 0x100000 /* drive can write MRW */
	CDC_RAM            = 0x200000 /* ok to open for WRITE */

	/* drive status possibilities returned by CDROM_DRIVE_STATUS ioctl */
	CDS_NO_INFO         = 0 /* if not implemented */
	CDS_NO_DISC         = 1
	CDS_TRAY_OPEN       = 2
	CDS_DRIVE_NOT_READY = 3
	CDS_DISC_OK         = 4

	/* return values for the CDROM_DISC_STATUS ioctl */
	/* can also return CDS_NO_[INFO|DISC], from above */
	CDS_AUDIO  = 100
	CDS_DATA_1 = 101
	CDS_DATA_2 = 102
	CDS_XA_2_1 = 103
	CDS_XA_2_2 = 104
	CDS_MIXED  = 105

	/* User-configurable behavior options for the uniform CD-ROM driver */
	CDO_AUTO_CLOSE = 0x1  /* close tray on first open() */
	CDO_AUTO_EJECT = 0x2  /* open tray on last release() */
	CDO_USE_FFLAGS = 0x4  /* use O_NONBLOCK information on open */
	CDO_LOCK       = 0x8  /* lock tray on open files */
	CDO_CHECK_TYPE = 0x10 /* check type on open for data */

	/* Special codes used when specifying changer slots. */
	CDSL_NONE    = math.MaxInt32 - 1
	CDSL_CURRENT = math.MaxInt32

	/* For partition based multisession access. IDE can handle 64 partitions
	 * per drive - SCSI CD-ROM's use minors to differentiate between the
	 * various drives, so we can't do multisessions the same way there.
	 * Use the -o session=x option to mount on them.
	 */
	CD_PART_MAX  = 64
	CD_PART_MASK = (CD_PART_MAX - 1)

	/*********************************************************************
	 * Generic Packet commands, MMC commands, and such
	 *********************************************************************/

	/* The generic packet command opcodes for CD/DVD Logical Units,
	 * From Table 57 of the SFF8090 Ver. 3 (Mt. Fuji) draft standard. */
	GPCMD_BLANK                         = 0xa1
	GPCMD_CLOSE_TRACK                   = 0x5b
	GPCMD_FLUSH_CACHE                   = 0x35
	GPCMD_FORMAT_UNIT                   = 0x04
	GPCMD_GET_CONFIGURATION             = 0x46
	GPCMD_GET_EVENT_STATUS_NOTIFICATION = 0x4a
	GPCMD_GET_PERFORMANCE               = 0xac
	GPCMD_INQUIRY                       = 0x12
	GPCMD_LOAD_UNLOAD                   = 0xa6
	GPCMD_MECHANISM_STATUS              = 0xbd
	GPCMD_MODE_SELECT_10                = 0x55
	GPCMD_MODE_SENSE_10                 = 0x5a
	GPCMD_PAUSE_RESUME                  = 0x4b
	GPCMD_PLAY_AUDIO_10                 = 0x45
	GPCMD_PLAY_AUDIO_MSF                = 0x47
	GPCMD_PLAY_AUDIO_TI                 = 0x48
	GPCMD_PLAY_CD                       = 0xbc
	GPCMD_PREVENT_ALLOW_MEDIUM_REMOVAL  = 0x1e
	GPCMD_READ_10                       = 0x28
	GPCMD_READ_12                       = 0xa8
	GPCMD_READ_BUFFER                   = 0x3c
	GPCMD_READ_BUFFER_CAPACITY          = 0x5c
	GPCMD_READ_CDVD_CAPACITY            = 0x25
	GPCMD_READ_CD                       = 0xbe
	GPCMD_READ_CD_MSF                   = 0xb9
	GPCMD_READ_DISC_INFO                = 0x51
	GPCMD_READ_DVD_STRUCTURE            = 0xad
	GPCMD_READ_FORMAT_CAPACITIES        = 0x23
	GPCMD_READ_HEADER                   = 0x44
	GPCMD_READ_TRACK_RZONE_INFO         = 0x52
	GPCMD_READ_SUBCHANNEL               = 0x42
	GPCMD_READ_TOC_PMA_ATIP             = 0x43
	GPCMD_REPAIR_RZONE_TRACK            = 0x58
	GPCMD_REPORT_KEY                    = 0xa4
	GPCMD_REQUEST_SENSE                 = 0x03
	GPCMD_RESERVE_RZONE_TRACK           = 0x53
	GPCMD_SEND_CUE_SHEET                = 0x5d
	GPCMD_SCAN                          = 0xba
	GPCMD_SEEK                          = 0x2b
	GPCMD_SEND_DVD_STRUCTURE            = 0xbf
	GPCMD_SEND_EVENT                    = 0xa2
	GPCMD_SEND_KEY                      = 0xa3
	GPCMD_SEND_OPC                      = 0x54
	GPCMD_SET_READ_AHEAD                = 0xa7
	GPCMD_SET_STREAMING                 = 0xb6
	GPCMD_START_STOP_UNIT               = 0x1b
	GPCMD_STOP_PLAY_SCAN                = 0x4e
	GPCMD_TEST_UNIT_READY               = 0x00
	GPCMD_VERIFY_10                     = 0x2f
	GPCMD_WRITE_10                      = 0x2a
	GPCMD_WRITE_12                      = 0xaa
	GPCMD_WRITE_AND_VERIFY_10           = 0x2e
	GPCMD_WRITE_BUFFER                  = 0x3b
	/* This is listed as optional in ATAPI 2.6, but is (curiously)
	 * missing from Mt. Fuji, Table 57.  It _is_ mentioned in Mt. Fuji
	 * Table 377 as an MMC command for SCSi devices though...  Most ATAPI
	 * drives support it. */
	GPCMD_SET_SPEED = 0xbb
	/* This seems to be a SCSI specific CD-ROM opcode
	 * to play data at track/index */
	GPCMD_PLAYAUDIO_TI = 0x48
	/*
	 * From MS Media Status Notification Support Specification. For
	 * older drives only.
	 */
	GPCMD_GET_MEDIA_STATUS = 0xda

	/* Mode page codes for mode sense/set */
	GPMODE_VENDOR_PAGE       = 0x00
	GPMODE_R_W_ERROR_PAGE    = 0x01
	GPMODE_WRITE_PARMS_PAGE  = 0x05
	GPMODE_WCACHING_PAGE     = 0x08
	GPMODE_AUDIO_CTL_PAGE    = 0x0e
	GPMODE_POWER_PAGE        = 0x1a
	GPMODE_FAULT_FAIL_PAGE   = 0x1c
	GPMODE_TO_PROTECT_PAGE   = 0x1d
	GPMODE_CAPABILITIES_PAGE = 0x2a
	GPMODE_ALL_PAGES         = 0x3f
	/* Not in Mt. Fuji, but in ATAPI 2.6 -- deprecated now in favor
	 * of MODE_SENSE_POWER_PAGE */
	GPMODE_CDROM_PAGE = 0x0d

	/* DVD struct types */
	DVD_STRUCT_PHYSICAL  = 0x00
	DVD_STRUCT_COPYRIGHT = 0x01
	DVD_STRUCT_DISCKEY   = 0x02
	DVD_STRUCT_BCA       = 0x03
	DVD_STRUCT_MANUFACT  = 0x04

	DVD_LAYERS = 4

	/*
	 * DVD authentication ioctl
	 */

	/* Authentication states */
	DVD_LU_SEND_AGID        = 0
	DVD_HOST_SEND_CHALLENGE = 1
	DVD_LU_SEND_KEY1        = 2
	DVD_LU_SEND_CHALLENGE   = 3
	DVD_HOST_SEND_KEY2      = 4

	/* Termination states */
	DVD_AUTH_ESTABLISHED = 5
	DVD_AUTH_FAILURE     = 6

	/* Other functions */
	DVD_LU_SEND_TITLE_KEY   = 7
	DVD_LU_SEND_ASF         = 8
	DVD_INVALIDATE_AGID     = 9
	DVD_LU_SEND_RPC_STATE   = 10
	DVD_HOST_SEND_RPC_STATE = 11

	DVD_CPM_NO_COPYRIGHT = 0
	DVD_CPM_COPYRIGHTED  = 1

	DVD_CP_SEC_NONE  = 0
	DVD_CP_SEC_EXIST = 1

	DVD_CGMS_UNRESTRICTED = 0
	DVD_CGMS_SINGLE       = 2
	DVD_CGMS_RESTRICTED   = 3

	/*
	 * feature profile
	 */
	CDF_RWRT = 0x0020 /* "Random Writable" */
	CDF_HWDM = 0x0024 /* "Hardware Defect Management" */
	CDF_MRW  = 0x0028

	/*
	 * media status bits
	 */
	CDM_MRW_NOTMRW            = 0
	CDM_MRW_BGFORMAT_INACTIVE = 1
	CDM_MRW_BGFORMAT_ACTIVE   = 2
	CDM_MRW_BGFORMAT_COMPLETE = 3

	/*
	 * mrw address spaces
	 */
	MRW_LBA_DMA = 0
	MRW_LBA_GAA = 1

	/*
	 * mrw mode pages (first is deprecated) -- probed at init time and
	 * cdi->mrw_mode_page is set
	 */
	MRW_MODE_PC_PRE1 = 0x2c
	MRW_MODE_PC      = 0x03
)
