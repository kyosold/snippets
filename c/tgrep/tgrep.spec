%define name tgrep
%define version 0.1
%define release 0
%define _prefix /usr/bin

# 关闭生成 debuginfo 包
%define debug_package %{nil}

Summary: find 'keyword' from big file with start to end hour
Name: %{name}
Version: %{version}
Release: %{release}
License: GPL 
Group: Development/Tools
Source: %{name}-%{version}.tar.gz
Distribution: RHEL6.5
Packager: kyosold@qq.com
AutoReqProv: no
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
find 'keyword' from big file with start to end hour

%prep
%setup -q

%build
[ "$RPM_BUILD_ROOT" != "/" ] && rm -rf $RPM_BUILD_ROOT
make

%install
echo "install " %{name}

mkdir -p %{buildroot}/%{_prefix}/

install -s -m 0755 ./%{name} %{buildroot}/%{_prefix}/%{name}

#安装后执行的脚本 
%post
echo "install finished."

%postun
#卸载包后执行下面的
rm -rf %{_prefix}/%{name}

%clean
rm -rf %{buildroot}

%files
%{_prefix}/%{name}

