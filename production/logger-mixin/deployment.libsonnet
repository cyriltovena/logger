{
  local container = $.core.v1.container,
  local deployment = $.apps.v1.deployment,

  logger_args:: {
    url: $._config.url,
    logps: $._config.logps,
    json: $._config.json,
  },

  logger_container::
    container.new('logger', $._images.logger) +
    container.withArgsMixin($.util.mapToFlags($.logger_args)) +
    $.util.resourcesRequests('500m', '10Mi') +
    $.util.resourcesLimits('1', '1Gi'),


  logger_deployment:
    deployment.new('logger', 10, [$.logger_container]),
}
