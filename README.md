# cert-exporter

A Prometheus exporter to expose certificate expiration dates for alerting.

If you're looking at adopting this, thank you. That said, you are
almost certainly better off looking at a process wherein you *don't*
have to monitor certificate expiration dates - think
[ACME](https://www.ietf.org/id/draft-ietf-acme-ari-01.html) and
[Let'sEncrypt](https://letsencrypt.org/).

## Usage

1. Create a config file in the style of the included `config.yaml`
   with as many host/port combinations as you care to monitor.
1. `cert-exporter -config /path/to/config.yaml`
1. Teach your prometheus instance to scrape the new target.
1. Write graphs and alerts (see caveats below).

## Caveats

I mean, it *should* be that simple, but there's always a surprise,
right?

### Multiply results by 1000 for Grafana

Grafana expects time/date data to be in microseconds(?) from
epoch. Precision concerns, I guess. Thus, to convert the values
exported here to actual dates, multiply the raw metric by 1000 first,
then treat as a date/time.

### We currently do not export intermediate cert data

Given the way the code (currently) loops through the certs that it
finds, for any given HostPort it probes, it will spit out the data for
the *last* cert that it finds in a chain. I do it this way because
that was the simplest possible thing, and I rationalize not doing
better because that's *probably* the one you're actually interested in
anyway.

## TODOs

1. Teach cert-exporter to export individually all the certs in a chain
   as separate metrics.
1. Teach cert-exporter to bind to alternative interfaces and ports.
1. Add a one-shot mode: i.e. immediately probe, print, and return -
   perhaps this would be useful to plug into other-than-prometheus
   monitoring pipelines.
1. Add a commandline or env-based config mode.
1. Extend the `config.yaml` to accommodate bind address/port settings.

## Motivation

This exporter exists for roughly two reasons:

1. (20+ years ago) I used to work at VeriSign, right next to the SSL
   Cert Support team, and
1. A more recent employer refused on a policy basis to (completely)
   automate cert renewal, at least internally.

Traditionally, when applying for an SSL cert, you'd supply a contact
email address, both for ID verification, as well as a renewal
contact. Most businesses getting certs (at least at the time) were
small, not highly technical, and not highly disciplined. The email
address provided would be an individual's, and very often, renewal
emails would be spam filtered, ignored, or even be sent to someone who
no longer worked at the company, and so the renewal process would
fail. You *could* ameliorate this by moving to a group mail address,
but you might still have filters or empty groups breaking the process.

Humans and manual processes *will* fail; we should just automate cert
renewal.

Enter a more recent employer, who went so far as to implement a full
internal CA with full ACME support. That should have been the end of
it, but out of an abundance of caution (certs being potentially useful
to attackers in multiple ways), said employer set a policy that
absolutely all certs issued from said CA had to be attributable to a
verified (2FA!)  internally employed human being. Thus, you could
script just about every aspect of creating a CSR, retrieving a
cert/key, building P12/keystore files, transporting them to machines,
updating configs, restarting daemons, whatever you liked - but there
still was required a physical human at a keyboard typing a password
and entering a 2FA token for every single renewal.

Welp. Monitoring it is.
