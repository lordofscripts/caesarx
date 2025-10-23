Name:           caesarx
Version:        1.1.1
Release:        1%{?dist}
Summary:        CaesarX is a collection of modernized Caesar ciphers for the real world.

License:        CC BY-NC-ND 4.0
URL:            https://github.com/lordofscripts/caesarx
Source:         %{name}-%{version}.tar.gz
BuildArch:      x86_64

Requires:       golang

%description
CaesarX is a collection of modernized Caesar ciphers (caesar, didimus, fibonacci, bellaso, vigenere, affine) supporting text and binary encryption and multiple modern-day alphabets.

%prep
%setup -q

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/%{_bindir}
install %{name} $RPM_BUILD_ROOT/%{_bindir}
mkdir -p $RPM_BUILD_ROOT/%{_sysconfdir}
install %{name}rc $RPM_BUILD_ROOT/%{_sysconfdir}
mkdir -p $RPM_BUILD_ROOT/%{_mandir}/man1/
install %{name}.1 $RPM_BUILD_ROOT/%{_mandir}/man1/

%post
ln -s -f $RPM_BUILD_ROOT/%{_bindir}/caesarx $RPM_BUILD_ROOT/%{_bindir}/didimus
ln -s -f $RPM_BUILD_ROOT/%{_bindir}/caesarx $RPM_BUILD_ROOT/%{_bindir}/fibonacci
ln -s -f $RPM_BUILD_ROOT/%{_bindir}/caesarx $RPM_BUILD_ROOT/%{_bindir}/bellaso
ln -s -f $RPM_BUILD_ROOT/%{_bindir}/caesarx $RPM_BUILD_ROOT/%{_bindir}/vigenere

%postun
rm -f $RPM_BUILD_ROOT/%{_bindir}/didimus
rm -f $RPM_BUILD_ROOT/%{_bindir}/fibonacci
rm -f $RPM_BUILD_ROOT/%{_bindir}/bellaso
rm -f $RPM_BUILD_ROOT/%{_bindir}/vigenere

%clean
rm -rf $RPM_BUILD_ROOT

%files
%{_bindir}/%{name}
%{_bindir}/affine
%{_bindir}/didimus -> %{_bindir}/%{name}
%{_bindir}/fibonacci -> %{_bindir}/%{name}
%{_bindir}/bellaso -> %{_bindir}/%{name}
%{_bindir}/vigenere -> %{_bindir}/%{name}
%{_bindir}/tabularecta
#%{_sysconfdir}/%{name}rc
%doc %{_mandir}/man1/%{name}.1.*
%doc %{_mandir}/man1/affine.1.*
%license LICENSE.md

%changelog
* Thu Oct 23 2025 lordofscripts
- configured spec file for use with GitHub Actions to automate building of RPM

* Sat Sep 27 2025 lordofscripts
- initial release