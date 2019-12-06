#!/usr/bin/env perl
use strict;
use warnings;
use List::Util qw(first);
use File::Slurp;
use YAML;

my $groups = Load(scalar(read_file("groups.yaml")));
my @lines = `git diff origin/master --name-only`;
my $code = $?;
if ($code != 0) {
    print "$_\n" for @lines;
    die("could not diff origin/master: exit code " . $code)
}

my @addons;
foreach my $line (@lines) {
    if ($line =~ m{^addons/([a-z]+)/?}) {
        push(@addons, $1);
    }
}

@addons = uniq(@addons);
foreach my $addon (@addons) {
    while(my($k, $v) = each %$groups) {
        if (first { $addon eq $_ } @$v) {
            my $group = "Test" . ucfirst($k) . "Group";
            my $exit = system("go test -race -v -run ${group}");
            if ($exit != 0) {
                die "tests failed for group ${group} (exit code ${code})";
            }
        }
    }
}

if (scalar(@addons) < 1) {
    exec "go test -run TestBuild"
}

# ------------------------------------------------------------------------------
# Private Functions
# ------------------------------------------------------------------------------

sub uniq {
    my %seen;
    grep !$seen{$_}++, @_;
}

1;
