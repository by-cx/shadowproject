package volumes

type VolumeInterface interface {
	Mount(source string, target string)
	Umount(target string)
}
