#!/bin/bash

cat >>/etc/my_init.d/services <<-EOGO
#!/bin/bash
dpkg-reconfigure openssh-server

/etc/init.d/postgresql start
/etc/init.d/redis-server start
/etc/init.d/ssh start
/usr/local/bin/MailHog &
EOGO
chmod a+x /etc/my_init.d/services

sed -i "s/#UsePAM/UsePAM/" /etc/ssh/sshd_config

exec /sbin/my_init