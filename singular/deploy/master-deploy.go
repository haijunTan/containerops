package deploy

import (
	cmd "github.com/Huawei/containerops/singular/cmd"
	"github.com/Huawei/containerops/singular/init_config"
)

func Deploymaster(list map[string]string, ip string) {

	opsfirewalld()

	echoconfig(list)
	opsetcd()
	opskube_apiserver()
	opskube_controller_manager()
	opskube_scheduler()
	opsflanneld()
	// #refresh service
	cmd.Reload()

	// #for SERVICES in etcd kube-apiserver kube-controller-manager kube-scheduler ; do
	// 	#systemctl restart $SERVICES
	// 	#systemctl enable $SERVICES
	// 	#systemctl status $SERVICES
	// #done
	cmd.RestartSvc([]string{"etcd", "kube-apiserver", "kube-controller-manager", "kube-scheduler"})

	opskubectlconfig()

	// // kubectl get nodes
	cmd.ExecCMDparams("kubectl", []string{"config", "cluster-info", "dump"})

}

func opsfirewalld() {

	// firewalld
	cmd.ExecCMDparams("systemctl", []string{"disable", "firewalld"})
	cmd.ExecCMDparams("systemctl", []string{"stop", "firewalld"})
}
func echoconfig(list map[string]string) {

	for k, v := range list {
		cmd.ExecShCommandEcho(k, v)
	}

}
func opsetcd() {

	//#etc create dir
	cmd.ExecCMDparams("mkdir", []string{"/etc/kubernetes/"})
	cmd.ExecCPparams("/tmp/config/config", "/etc/kubernetes/config") //
	// #etc service
	cmd.ExecCMDparams("mkdir", []string{"/etc/etcd/"})
	cmd.ExecCPparams("/tmp/etcd.conf", "/etc/etcd/etcd.conf")

	cmd.ExecCMDparams("mkdir", []string{"/usr/lib/systemd/system/"})
	cmd.ExecCPparams("/tmp/config/etcd.service", "/usr/lib/systemd/system/etcd.service") //
	cmd.ExecCPparams("/tmp/etcd", "/usr/bin/etcd")
	cmd.ExecCPparams("/tmp/etcdctl", "/usr/bin/etcdctl")
	cmd.ExecCMDparams("mkdir", []string{"/var/lib/etcd/"})
	// #/var/lib/etcd/default.etcd auto?

	cmd.ExecCMDparams("systemctl", []string{"daemon-reload"})
	cmd.ServiceStart("etcd")
	cmd.ServiceIsEnabled("etcd")
	cmd.ServiceExists("etcd")
	cmd.ExecCMDparams("etcdctl", []string{"mkdir", init_config.EtcdNet}) //"/kube-centos/network/config nend update config
	cmd.ExecCMDparams("etcdctl", []string{"mk", init_config.EtcdNet, "{\"Network\":\"172.40.0.0/16\"\\,\"SubnetLen\":24\\,\"Backend\":{\"Type\":\"vxlan\"}}"})

}

func opskube_apiserver() {
	//#kube-apiserver
	cmd.ExecCPparams("/tmp/kube-apiserver", "/usr/bin/kube-apiserver")

	cmd.ExecCPparams("/tmp/config/kube-apiserver.service", "/usr/lib/systemd/system/kube-apiserver.service")
	cmd.ExecCPparams("/tmp/config/apiserver", "/etc/kubernetes/apiserver")

}
func opskube_controller_manager() {
	//#kube-controller-manager
	cmd.ExecCPparams("/tmp/kube-controller-manager", "/usr/bin/kube-controller-manager")

	cmd.ExecCPparams("/tmp/config/kube-controller-manager.service", "/usr/lib/systemd/system/kube-controller-manager.service")
	cmd.ExecCPparams("/tmp/config/controller-manager", "/etc/kubernetes/controller-manager")
}
func opskube_scheduler() {

	//#kube-scheduler
	cmd.ExecCPparams("/tmp/kube-scheduler", "/usr/bin/kube-scheduler")

	cmd.ExecCPparams("/tmp/config/kube-scheduler.service", "/usr/lib/systemd/system/kube-scheduler.service")
	cmd.ExecCPparams("/tmp/config/scheduler", "/etc/kubernetes/scheduler")

}
func opsflanneld() {
	// #flanneld reserved
	cmd.ExecCPparams("/tmp/flanneld", "/usr/bin/flanneld")

	cmd.ExecCMDparams("mkdir", []string{"-p", "/usr/libexec/flannel/"})
	cmd.ExecCMDparams("mkdir", []string{"/tmp/etcd"})
	// check work
	cmd.ExecCPparams("/tmp/flannel/usr/libexec/flannel/mk-docker-opts.sh", "/usr/libexec/flannel/mk-docker-opts.sh")
	cmd.ExecCMDparams("mkdir", []string{"-p", "/etc/sysconfig/"})

	cmd.ExecCPparams("/tmp/config/flanneld", "/etc/sysconfig/flanneld")

	cmd.ExecCPparams("/tmp/config/flanneld.service", "/usr/lib/systemd/system/flanneld.service")

}

func opskubectlconfig() {
	// #kubectl config
	cmd.ExecCPparams("/tmp/kubectl", "/usr/bin/kubectl")

	cmd.ExecCMDparams("kubectl", []string{"config", "set-cluster", "default-cluster", "--server=http://centos-master:8080"})
	cmd.ExecCMDparams("kubectl", []string{"config", "set-cluster", "default-cluster", "--cluster=default-cluster", "--user=default-admin"})
	cmd.ExecCMDparams("kubectl", []string{"config", "use-context", "use-context"})
	cmd.ExecCMDparams("kubectl", []string{"config", "get", "nodes"})

}
