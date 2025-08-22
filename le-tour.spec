Name:           le-tour-dashmore-server
Version:        %{_version}
Release:        1%{?dist}
Summary:        The API server for Le Tour d'Ashmore

License:        GPL
URL:            https://rpm.joshuacrunden.com
Source0:        file://%{name}-%{_version}.tar.gz

Requires: postgresql-server

%define appuser letour
%global debug_package %{nil}
%global source_date_epoch_from_changelog %{nil}
%define app_config_dir %{_sysconfdir}/%{name}
%define config_file server.conf
%define dbsetup database-setup.sql
%define setup_script le-tour-setup.sh
%define setup_service le-tour-setup.service

%description
The API server for Le Tour d'Ashmore

%prep
%setup -q

%build

%check

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/%{_bindir}
mkdir -p $RPM_BUILD_ROOT/%{_unitdir}
mkdir -p $RPM_BUILD_ROOT/%{app_config_dir}

cp %{name} $RPM_BUILD_ROOT/%{_bindir}
cp %{name}.service $RPM_BUILD_ROOT/%{_unitdir}
cp %{setup_service} $RPM_BUILD_ROOT/%{_unitdir}
cp %{config_file} $RPM_BUILD_ROOT/%{app_config_dir}
cp %{dbsetup} $RPM_BUILD_ROOT/%{app_config_dir}
cp %{setup_script} $RPM_BUILD_ROOT/%{_bindir}

%pre
getent passwd %{appuser} > /dev/null
if [ $? -ne 0 ]
then
    useradd %{appuser}
fi
# initialise postgresql database if it hasn't already been done
if [ ! -d "/var/lib/pgsql/data" ]
then
        /usr/bin/postgresql-setup --initdb
fi

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service
getent passwd %{appuser} > /dev/null
# if the user exists and it's an uninstall (not upgrade)
if [ $? -eq 0 ] && [ $1 -eq 0 ]
then
    userdel %{appuser}
    rm -rf /home/letour
    rm -rf /var/spool/mail/letour
fi

%postun
%systemd_postun_with_restart %{name}.service

%files
%{_bindir}/%{name}
%{_unitdir}/%{name}.service
%{_unitdir}/%{setup_service}
%{app_config_dir}/%{config_file}
%{app_config_dir}/%{dbsetup}
%{_bindir}/%{setup_script}
%dir %{app_config_dir}
