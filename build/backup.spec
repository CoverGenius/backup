# SPEC file overview:
# https://docs.fedoraproject.org/en-US/quick-docs/creating-rpm-packages/#con_rpm-spec-file-overview
# Fedora packaging guidelines:
# https://docs.fedoraproject.org/en-US/packaging-guidelines/


Name: backup
Version: 0.0.8
Release: 0%{?dist}
Summary: Backup tool
License: BSD
URL: https://www.covergenius.com
BuildRequires: golang >= 1.12
BuildArch: x86_64


%description
Backup tool


%prep
mkdir -p %{_topdir}/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
cp -rf %{_sourcedir}/* %{_topdir}/BUILD


%build
go build


%install
rm -rf %{buildroot}
%{__install} -D -m 0644 %{_topdir}/BUILD/backup %{buildroot}/%{_sbindir}/backup


%files
%defattr(755,root,root,755)
%{_sbindir}/backup


%changelog
* Wed Sep 23 2020 Serghei Anicheev <serghei@covergenius.com>
- Added a flag which does not delete tmp directory
* Tue May 19 2020 Serghei Anicheev <serghei@covergenius.com>
- Added MySQL source
- Added restore logic
* Thu Jan 22 2020 Serghei Anicheev <serghei@covergenius.com>
- Adding db_snapshot_suffix option
* Thu Jan 2 2020 Serghei Anicheev <serghei@covergenius.com>
- Initial commit
