# If /etc/vbox doesn't exist create it
if [ ! -d /etc/vbox ]; then
    echo "Creating directory /etc/vbox"
    mkdir /etc/vbox
fi

# if file /etc/vbox/networks.conf doesn't exist create it
if [ ! -f /etc/vbox/networks.conf ]; then
    echo "Creating file /etc/vbox/networks.conf"
    touch /etc/vbox/networks.conf
fi

# if file /etc/vbox/networks.conf doesn't contain "* 172.28.128.0/16" add it
if ! grep -q "* 172.28.128.0/16" /etc/vbox/networks.conf; then
    echo "Adding '* 172.28.128.0/16' to /etc/vbox/networks.conf"
    echo "* 172.28.128.0/16" > /etc/vbox/networks.conf
fi

# if file /etc/vbox/networks.conf doesn't contain "* 172.28.128.0/16" add it
if ! grep -q "* 172.28.128.0/16" /etc/vbox/networks.conf; then
    echo "Adding '* 172.28.128.0/16' to /etc/vbox/networks.conf"
    echo "* 172.28.128.0/16" >> /etc/vbox/networks.conf
fi

# if file /etc/vbox/networks.conf doesn't contain "* 192.168.20.0/16" add it
# This is necessary for the client network
if ! grep -q "* 192.168.20.0/16" /etc/vbox/networks.conf; then
    echo "Adding '* 192.168.20.0/16' to /etc/vbox/networks.conf"
    echo "* 192.168.20.0/16" >> /etc/vbox/networks.conf
fi