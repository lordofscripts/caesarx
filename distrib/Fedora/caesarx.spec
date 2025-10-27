%global debug_package %{nil} %global _builddebuginfo_packages %{nil}
%define bindir /usr/local/bin
%define __cp %{__cp} -r

Name:           caesarx
Version:        1.1.1
Release:        1%{?dist}
Summary:        CaesarX is a collection of modernized Caesar ciphers for the real world.

License:        CC BY-NC-ND 4.0
URL:            https://github.com/lordofscripts/caesarx
Source:         %{name}-%{version}.tar.gz
BuildArch:      x86_64
BuildRequires:  golang
Requires:       golang

%description
CaesarX is a collection of modernized Caesar ciphers (caesar, didimus, fibonacci, bellaso, vigenere, affine) supporting text and binary encryption and multiple modern-day alphabets.

%prep
%setup -q

%build
[ -d bin ] || mkdir bin
go build -tags logx -v -buildmode=pie -o bin/caesarx cmd/caesar/*.go
go build -tags logx -v -buildmode=pie -o bin/affine cmd/affine/*go
go build -tags logx -v -buildmode=pie -o bin/tabularecta cmd/tabularecta/*go

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/%{_bindir}
install -m 0755  bin/%{name} $RPM_BUILD_ROOT/%{_bindir}
install -m 0755  bin/affine $RPM_BUILD_ROOT/%{_bindir}
install -m 0755  bin/tabularecta $RPM_BUILD_ROOT/%{_bindir}
ln -s %{_bindir}/%{name} %{buildroot}%{_bindir}/didimus
ln -s %{_bindir}/%{name} %{buildroot}%{_bindir}/fibonacci
ln -s %{_bindir}/%{name} %{buildroot}%{_bindir}/bellaso
ln -s %{_bindir}/%{name} %{buildroot}%{_bindir}/vigenere
#mkdir -p $RPM_BUILD_ROOT/%{_sysconfdir}
#install %{name}rc $RPM_BUILD_ROOT/%{_sysconfdir}
mkdir -p $RPM_BUILD_ROOT/%{_mandir}/man1/
install distrib/manpages/man1/%{name}.1 $RPM_BUILD_ROOT/%{_mandir}/man1/
install distrib/manpages/man1/affine.1 $RPM_BUILD_ROOT/%{_mandir}/man1/
install distrib/manpages/man1/didimus.1 $RPM_BUILD_ROOT/%{_mandir}/man1/
install distrib/manpages/man1/fibonacci.1 $RPM_BUILD_ROOT/%{_mandir}/man1/
install distrib/manpages/man1/bellaso.1 $RPM_BUILD_ROOT/%{_mandir}/man1/
install distrib/manpages/man1/vigenere.1 $RPM_BUILD_ROOT/%{_mandir}/man1/

%postun
rm -f $RPM_BUILD_ROOT/%{_bindir}/caesarx
rm -f $RPM_BUILD_ROOT/%{_bindir}/didimus
rm -f $RPM_BUILD_ROOT/%{_bindir}/fibonacci
rm -f $RPM_BUILD_ROOT/%{_bindir}/bellaso
rm -f $RPM_BUILD_ROOT/%{_bindir}/vigenere
rm -f $RPM_BUILD_ROOT/%{_bindir}/affine
rm -f $RPM_BUILD_ROOT/%{_bindir}/tabularecta
rm -f $RPM_BUILD_ROOT/%{_mandir}/man1/%{name}.1.gz
rm -f $RPM_BUILD_ROOT/%{_mandir}/man1/affine.1.gz
rm -f $RPM_BUILD_ROOT/%{_mandir}/man1/bellaso.1.gz
rm -f $RPM_BUILD_ROOT/%{_mandir}/man1/didimus.1.gz
rm -f $RPM_BUILD_ROOT/%{_mandir}/man1/fibonacci.1.gz
rm -f $RPM_BUILD_ROOT/%{_mandir}/man1/vigenere.1.gz

%clean
rm -rf $RPM_BUILD_ROOT

%files
%{_bindir}/%{name}
%{_bindir}/affine
%{_bindir}/tabularecta
%{_bindir}/didimus
%{_bindir}/fibonacci
%{_bindir}/bellaso
%{_bindir}/vigenere
%{_mandir}/man1/%{name}.1.gz
%{_mandir}/man1/affine.1.gz
%{_mandir}/man1/bellaso.1.gz
%{_mandir}/man1/didimus.1.gz
%{_mandir}/man1/fibonacci.1.gz
%{_mandir}/man1/vigenere.1.gz
#%{_sysconfdir}/%{name}rc
%doc %{_mandir}/man1/%{name}.1.*
%doc %{_mandir}/man1/affine.1.*
%license LICENSE.md

%changelog
* Mon Oct 27 2025 lordofscripts
- added missing manual pages to the install

* Thu Oct 23 2025 lordofscripts
- configured spec file for use with GitHub Actions to automate building of RPM

* Sat Sep 27 2025 lordofscripts
- initial release