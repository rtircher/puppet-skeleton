class tircher::base {

#  if ! defined(Package['git']) { package { 'git': ensure => installed } }

}
