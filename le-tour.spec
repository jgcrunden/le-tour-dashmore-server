Name:           le-tour-dashmore-server
Version:        %{_version}
Release:        1%{?dist}
Summary:        The API server for Le Tour d'Ashmore

License:        GPL
URL:            https://rpm.joshuacrunden.com
Source0:        file://%{name}-%{_version}.tar.gz

Requires: postgresql

%define appuser letour
%global debug_package %{nil}

%description
The API server for Le Tour d'Ashmore

%prep
%setup -q

%build

%check

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/%{_bindir}
pwd
cp %{name} $RPM_BUILD_ROOT/%{_bindir}

%pre
getent passwd %{appuser}
if [ $? -ne 0 ]
then
    useradd %{appuser}
fi

%preun
getent passwd %{appuser}
if [ $? -eq 0 ]
then
    userdel %{appuser}
fi



%files
%{_bindir}/%{name}

